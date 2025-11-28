package document

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/pkg/gotenberg"
	"techmind/schema/ent"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	// MaxFileSize - максимальный размер файла 5 ГБ
	MaxFileSize = 5 * 1024 * 1024 * 1024 // 5GB в байтах
)

// AllowedMimeTypes - список разрешенных MIME типов
var AllowedMimeTypes = map[string]bool{
	// Документы
	"application/pdf":    true,
	"application/msword": true, // .doc
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true, // .docx
	"application/vnd.ms-excel": true, // .xls
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true, // .xlsx
	"application/vnd.ms-powerpoint":                                             true, // .ppt
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true, // .pptx
	"text/plain":      true, // .txt
	"text/csv":        true, // .csv
	"application/rtf": true, // .rtf

	// Изображения
	"image/jpeg":    true, // .jpg, .jpeg
	"image/png":     true, // .png
	"image/gif":     true, // .gif
	"image/webp":    true, // .webp
	"image/svg+xml": true, // .svg
	"image/bmp":     true, // .bmp
	"image/tiff":    true, // .tiff

	// Видео (потенциально)
	"video/mp4":        true, // .mp4
	"video/mpeg":       true, // .mpeg
	"video/quicktime":  true, // .mov
	"video/x-msvideo":  true, // .avi
	"video/x-matroska": true, // .mkv
	"video/webm":       true, // .webm
}

// AllowedExtensions - список разрешенных расширений файлов
var AllowedExtensions = map[string]bool{
	// Документы
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".txt":  true,
	".csv":  true,
	".rtf":  true,

	// Изображения
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
	".bmp":  true,
	".tiff": true,
	".tif":  true,

	// Видео
	".mp4":  true,
	".mpeg": true,
	".mpg":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".webm": true,
}

type documentService struct {
	documentRepo    repo.DocumentRepository
	documentTagRepo repo.DocumentTagRepository
	tagRepo         repo.TagRepository
	folderRepo      repo.FolderRepository
	minioClient     *minio.Client
	bucketName      string
	gotenbergClient *gotenberg.Client
}

func NewService(
	documentRepo repo.DocumentRepository,
	documentTagRepo repo.DocumentTagRepository,
	tagRepo repo.TagRepository,
	folderRepo repo.FolderRepository,
	minioClient *minio.Client,
	gotenbergClient *gotenberg.Client,
) service.DocumentService {

	return &documentService{
		documentRepo:    documentRepo,
		documentTagRepo: documentTagRepo,
		tagRepo:         tagRepo,
		folderRepo:      folderRepo,
		minioClient:     minioClient,
		bucketName:      "documents",
		gotenbergClient: gotenbergClient,
	}
}

func (s *documentService) Upload(ctx context.Context, input service.DocumentUploadInput) (*ent.Document, error) {
	// Проверяем размер файла
	if input.FileSize > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of 5GB")
	}

	// Проверяем расширение файла
	ext := strings.ToLower(filepath.Ext(input.Name))
	if !AllowedExtensions[ext] {
		return nil, fmt.Errorf("file type not supported: %s", ext)
	}

	// Проверяем MIME тип
	if !AllowedMimeTypes[input.MimeType] {
		return nil, fmt.Errorf("file type not supported: %s", input.MimeType)
	}

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
	objectName := fmt.Sprintf("%s/%s%s", input.CompanyID.String(), fileID.String(), ext)

	// ...existing code...

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

		// Проверяем на конфликт уникальности checksum
		if ent.IsConstraintError(err) {
			return nil, fmt.Errorf("document with this checksum already exists")
		}

		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	// Генерация preview для поддерживаемых типов файлов
	if s.isConvertibleToPDF(input.MimeType) {
		// Запускаем генерацию preview асинхронно, чтобы не блокировать загрузку
		go func() {
			// Создаем новый контекст с таймаутом для фоновой задачи
			previewCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			if err := s.GeneratePDFPreview(previewCtx, document.ID); err != nil {
				// Логируем ошибку, но не прерываем процесс загрузки
				// В продакшене здесь должно быть логирование через logger
				fmt.Printf("Failed to generate preview for document %s: %v\n", document.ID, err)
			}
		}()
	}

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

	// Получаем download URL
	downloadURL, _ := s.GetDownloadURL(ctx, documentID)

	return &service.DocumentWithTags{
		Document:    document,
		Tags:        tags,
		PreviewURL:  previewURL,
		DownloadURL: downloadURL,
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
		downloadURL, _ := s.GetDownloadURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:    doc,
			Tags:        tags,
			PreviewURL:  previewURL,
			DownloadURL: downloadURL,
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
		downloadURL, _ := s.GetDownloadURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:    doc,
			Tags:        tags,
			PreviewURL:  previewURL,
			DownloadURL: downloadURL,
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
		downloadURL, _ := s.GetDownloadURL(ctx, doc.ID)

		result = append(result, &service.DocumentWithTags{
			Document:    doc,
			Tags:        tags,
			PreviewURL:  previewURL,
			DownloadURL: downloadURL,
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

// GeneratePDFPreview конвертирует файл документа в PDF превью и загружает его в MinIO
// Поддерживает конвертацию Office документов (docx, xlsx, pptx и т.д.) через Gotenberg
// После успешной конвертации обновляет ссылку на preview в базе данных
func (s *documentService) GeneratePDFPreview(ctx context.Context, documentID uuid.UUID) error {
	// Проверяем что Gotenberg доступен
	//if !s.gotenbergEnabled || s.gotenbergClient == nil {
	//	return fmt.Errorf("gotenberg is not enabled or configured")
	//}

	// Получаем документ из БД
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Проверяем что документ поддерживает конвертацию
	if !s.isConvertibleToPDF(document.MimeType) {
		return fmt.Errorf("document type %s is not convertible to PDF", document.MimeType)
	}

	// Скачиваем оригинальный файл из MinIO
	object, err := s.minioClient.GetObject(ctx, s.bucketName, document.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get file from minio: %w", err)
	}
	defer object.Close()

	// Читаем содержимое файла в память
	fileContent, err := io.ReadAll(object)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Определяем имя файла для конвертации
	fileName := filepath.Base(document.FilePath)

	// Конвертируем в PDF через Gotenberg
	var pdfResponse *gotenberg.Response

	// Используем LibreOffice для Office документов
	if s.isOfficeDocument(document.MimeType) {
		pdfResponse, err = s.gotenbergClient.ConvertOfficeToPDF(
			ctx,
			[]gotenberg.File{
				{
					Name:    fileName,
					Content: fileContent,
				},
			},
			&gotenberg.LibreOfficeRequest{
				Landscape:        false,
				SinglePageSheets: true,
				OutputFilename:   "preview",
			},
		)
		if err != nil {
			return fmt.Errorf("failed to convert office document to PDF: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported document type for conversion: %s", document.MimeType)
	}

	// Генерируем путь для preview файла
	previewID := uuid.New()
	previewObjectName := fmt.Sprintf(
		"%s/previews/%s.pdf",
		document.CompanyID.String(),
		previewID.String(),
	)

	// Загружаем PDF preview в MinIO
	_, err = s.minioClient.PutObject(
		ctx,
		s.bucketName,
		previewObjectName,
		bytes.NewReader(pdfResponse.Body),
		int64(len(pdfResponse.Body)),
		minio.PutObjectOptions{
			ContentType: "application/pdf",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to upload preview to minio: %w", err)
	}

	// Обновляем путь к preview в базе данных
	err = s.documentRepo.UpdatePreviewPath(ctx, documentID, previewObjectName)
	if err != nil {
		// Если не удалось обновить БД, удаляем загруженный preview
		_ = s.minioClient.RemoveObject(ctx, s.bucketName, previewObjectName, minio.RemoveObjectOptions{})
		return fmt.Errorf("failed to update preview path in database: %w", err)
	}

	return nil
}

// isConvertibleToPDF проверяет, можно ли сконвертировать документ в PDF
func (s *documentService) isConvertibleToPDF(mimeType string) bool {
	convertibleTypes := []string{
		// Microsoft Office
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",   // docx
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",         // xlsx
		"application/vnd.openxmlformats-officedocument.presentationml.presentation", // pptx
		"application/msword",            // doc
		"application/vnd.ms-excel",      // xls
		"application/vnd.ms-powerpoint", // ppt

		// OpenDocument
		"application/vnd.oasis.opendocument.text",         // odt
		"application/vnd.oasis.opendocument.spreadsheet",  // ods
		"application/vnd.oasis.opendocument.presentation", // odp

		// Rich Text
		"application/rtf",
		"text/rtf",

		// HTML
		"text/html",
	}

	for _, ct := range convertibleTypes {
		if strings.EqualFold(mimeType, ct) {
			return true
		}
	}

	return false
}

// isOfficeDocument проверяет, является ли документ Office файлом
func (s *documentService) isOfficeDocument(mimeType string) bool {
	officeTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/msword",
		"application/vnd.ms-excel",
		"application/vnd.ms-powerpoint",
		"application/vnd.oasis.opendocument.text",
		"application/vnd.oasis.opendocument.spreadsheet",
		"application/vnd.oasis.opendocument.presentation",
		"application/rtf",
		"text/rtf",
	}

	for _, ot := range officeTypes {
		if strings.EqualFold(mimeType, ot) {
			return true
		}
	}

	return false
}
