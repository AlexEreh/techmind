package document

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type documentService struct {
	documentRepo    repo.DocumentRepository
	documentTagRepo repo.DocumentTagRepository
	tagRepo         repo.TagRepository
	folderRepo      repo.FolderRepository
	minioClient     *minio.Client
	bucketName      string
}

func NewService(
	documentRepo repo.DocumentRepository,
	documentTagRepo repo.DocumentTagRepository,
	tagRepo repo.TagRepository,
	folderRepo repo.FolderRepository,
	minioClient *minio.Client,
) service.DocumentService {
	return &documentService{
		documentRepo:    documentRepo,
		documentTagRepo: documentTagRepo,
		tagRepo:         tagRepo,
		folderRepo:      folderRepo,
		minioClient:     minioClient,
		bucketName:      "documents",
	}
}

func (s *documentService) Upload(ctx context.Context, input service.DocumentUploadInput) (*ent.Document, error) {
	// Проверяем что папка существует и принадлежит компании
	if input.FolderID != nil {
		folder, err := s.folderRepo.GetByID(ctx, *input.FolderID)
		if err != nil {
			return nil, fmt.Errorf("folder not found: %w", err)
		}
		if folder.CompanyID != input.CompanyID {
			return nil, fmt.Errorf("folder belongs to different company")
		}
	}

	// Генерируем уникальное имя файла
	fileID := uuid.New()
	ext := filepath.Ext(input.Name)
	objectName := fmt.Sprintf("%s/%s%s", input.CompanyID.String(), fileID.String(), ext)

	// Вычисляем checksum
	hash := sha256.New()
	teeReader := io.TeeReader(input.File, hash)

	// Загружаем файл в MinIO
	_, err := s.minioClient.PutObject(ctx, s.bucketName, objectName, teeReader, input.FileSize, minio.PutObjectOptions{
		ContentType: input.MimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to minio: %w", err)
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))

	// Создаем запись в БД
	document, err := s.documentRepo.Create(
		ctx,
		input.CompanyID,
		input.FolderID,
		input.Name,
		objectName,
		input.FileSize,
		input.MimeType,
		checksum,
	)
	if err != nil {
		// Удаляем файл из MinIO если не удалось создать запись в БД
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	// TODO: Генерация preview для поддерживаемых типов файлов скорее всего в микросервисе (kafka - посредник, без grpc)

	return document, nil
}

func (s *documentService) GetByID(ctx context.Context, documentID uuid.UUID) (*service.DocumentWithTags, error) {
	// Получаем документ
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Получаем теги документа
	tags, err := s.getDocumentTags(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document tags: %w", err)
	}

	// Получаем preview URL
	previewURL, _ := s.GetPreviewURL(ctx, documentID)

	return &service.DocumentWithTags{
		Document:   document,
		Tags:       tags,
		PreviewURL: previewURL,
	}, nil
}

func (s *documentService) GetByFolder(ctx context.Context, folderID uuid.UUID) ([]*service.DocumentWithTags, error) {
	// Получаем документы в папке
	documents, err := s.documentRepo.ListByFolder(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by folder: %w", err)
	}

	// Обогащаем документы тегами и preview URLs
	result := make([]*service.DocumentWithTags, 0, len(documents))
	for _, doc := range documents {
		tags, _ := s.getDocumentTags(ctx, doc.ID)
		previewURL, _ := s.GetPreviewURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:   doc,
			Tags:       tags,
			PreviewURL: previewURL,
		})
	}

	return result, nil
}

func (s *documentService) GetByCompany(ctx context.Context, companyID uuid.UUID) ([]*service.DocumentWithTags, error) {
	// Получаем все документы компании
	documents, err := s.documentRepo.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by company: %w", err)
	}

	// Обогащаем документы тегами и preview URLs
	result := make([]*service.DocumentWithTags, 0, len(documents))
	for _, doc := range documents {
		tags, _ := s.getDocumentTags(ctx, doc.ID)
		previewURL, _ := s.GetPreviewURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:   doc,
			Tags:       tags,
			PreviewURL: previewURL,
		})
	}

	return result, nil
}

func (s *documentService) Update(ctx context.Context, documentID uuid.UUID, input service.DocumentUpdateInput) (*ent.Document, error) {
	// Получаем документ
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Если меняется папка, проверяем что она существует и принадлежит той же компании
	if input.FolderID != nil {
		folder, err := s.folderRepo.GetByID(ctx, *input.FolderID)
		if err != nil {
			return nil, fmt.Errorf("folder not found: %w", err)
		}
		if folder.CompanyID != document.CompanyID {
			return nil, fmt.Errorf("folder belongs to different company")
		}
	}

	updatedDocument, err := s.documentRepo.Update(
		ctx,
		documentID,
		input.FolderID,
		input.SenderID,
		input.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	return updatedDocument, nil
}

func (s *documentService) Delete(ctx context.Context, documentID uuid.UUID) error {
	// Получаем документ
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Удаляем файлы из MinIO
	if err := s.minioClient.RemoveObject(ctx, s.bucketName, document.FilePath, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}

	// Удаляем preview если есть
	if document.PreviewFilePath != nil {
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, *document.PreviewFilePath, minio.RemoveObjectOptions{})
	}

	// Удаляем запись из БД (каскадно удалятся связи с тегами)
	if err := s.documentRepo.Delete(ctx, documentID); err != nil {
		return fmt.Errorf("failed to delete document record: %w", err)
	}

	return nil
}

func (s *documentService) GetDownloadURL(ctx context.Context, documentID uuid.UUID) (string, error) {
	// Получаем документ
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return "", fmt.Errorf("document not found: %w", err)
	}

	// Генерируем presigned URL на 1 час
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, document.FilePath, 1*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate download url: %w", err)
	}

	return url.String(), nil
}

func (s *documentService) GetPreviewURL(ctx context.Context, documentID uuid.UUID) (string, error) {
	// Получаем документ
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return "", fmt.Errorf("document not found: %w", err)
	}

	// Если нет preview, возвращаем пустую строку
	if document.PreviewFilePath == nil {
		return "", nil
	}

	// Генерируем presigned URL на 1 час
	url, err := s.minioClient.PresignedGetObject(ctx, s.bucketName, *document.PreviewFilePath, 1*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate preview url: %w", err)
	}

	return url.String(), nil
}

func (s *documentService) Search(ctx context.Context, companyID uuid.UUID, query string, folderID *uuid.UUID, tagIDs []uuid.UUID) ([]*service.DocumentWithTags, error) {
	// TODO: Реализовать поиск через Elasticsearch
	// Пока используем простую фильтрацию через БД

	var documents []*ent.Document
	var err error

	if folderID != nil {
		documents, err = s.documentRepo.ListByFolder(ctx, *folderID)
	} else {
		documents, err = s.documentRepo.ListByCompany(ctx, companyID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	// Фильтруем по тегам если указаны
	if len(tagIDs) > 0 {
		var filtered []*ent.Document
		for _, doc := range documents {
			docTags, _ := s.getDocumentTags(ctx, doc.ID)
			if s.hasAllTags(docTags, tagIDs) {
				filtered = append(filtered, doc)
			}
		}
		documents = filtered
	}

	// Обогащаем результаты
	result := make([]*service.DocumentWithTags, 0, len(documents))
	for _, doc := range documents {
		tags, _ := s.getDocumentTags(ctx, doc.ID)
		previewURL, _ := s.GetPreviewURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:   doc,
			Tags:       tags,
			PreviewURL: previewURL,
		})
	}

	return result, nil
}

func (s *documentService) getDocumentTags(ctx context.Context, documentID uuid.UUID) ([]*ent.Tag, error) {
	// Получаем связи документ-тег
	docTags, err := s.documentTagRepo.ListByDocument(ctx, documentID)
	if err != nil {
		return []*ent.Tag{}, nil
	}

	// Получаем полные объекты тегов
	tags := make([]*ent.Tag, 0, len(docTags))
	for _, dt := range docTags {
		tag, err := s.tagRepo.GetByID(ctx, dt.TagID)
		if err == nil {
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func (s *documentService) hasAllTags(docTags []*ent.Tag, tagIDs []uuid.UUID) bool {
	tagMap := make(map[uuid.UUID]bool)
	for _, tag := range docTags {
		tagMap[tag.ID] = true
	}

	for _, tagID := range tagIDs {
		if !tagMap[tagID] {
			return false
		}
	}

	return true
}
