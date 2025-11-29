package document

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/pkg/gotenberg"
	"techmind/schema/ent"

	"code.sajari.com/docconv/v2/client"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
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
	documentRepo        repo.DocumentRepository
	documentTagRepo     repo.DocumentTagRepository
	tagRepo             repo.TagRepository
	folderRepo          repo.FolderRepository
	minioClient         *minio.Client
	bucketName          string
	gotenbergClient     *gotenberg.Client
	elasticsearchClient *elasticsearch.Client
}

func NewService(
	documentRepo repo.DocumentRepository,
	documentTagRepo repo.DocumentTagRepository,
	tagRepo repo.TagRepository,
	folderRepo repo.FolderRepository,
	minioClient *minio.Client,
	gotenbergClient *gotenberg.Client,
	elasticsearchClient *elasticsearch.Client,
) service.DocumentService {

	return &documentService{
		documentRepo:        documentRepo,
		documentTagRepo:     documentTagRepo,
		tagRepo:             tagRepo,
		folderRepo:          folderRepo,
		minioClient:         minioClient,
		bucketName:          "documents",
		gotenbergClient:     gotenbergClient,
		elasticsearchClient: elasticsearchClient,
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

	// Извлечение текста и индексация в Elasticsearch
	if s.isExtractableText(input.MimeType) {
		// Запускаем извлечение текста асинхронно
		go func() {
			// Создаем новый контекст с таймаутом для фоновой задачи
			extractCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			if err := s.ExtractAndIndexText(extractCtx, document.ID); err != nil {
				// Логируем ошибку, но не прерываем процесс загрузки
				fmt.Printf("Failed to extract and index text for document %s: %v\n", document.ID, err)
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
	var documents []*ent.Document
	var err error

	// Если есть поисковый запрос, используем ТОЛЬКО Elasticsearch для полнотекстового поиска
	if query != "" {
		if s.elasticsearchClient == nil {
			return nil, fmt.Errorf("elasticsearch is not configured")
		}

		documentIDs, err := s.searchInElasticsearch(ctx, companyID, query, folderID)
		if err != nil {
			return nil, fmt.Errorf("elasticsearch search failed: %w", err)
		}

		// Если ничего не найдено в Elasticsearch, возвращаем пустой результат
		if len(documentIDs) == 0 {
			return []*service.DocumentWithTags{}, nil
		}

		// Получаем документы по ID из Elasticsearch
		documents, err = s.getDocumentsByIDs(ctx, documentIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get documents: %w", err)
		}
	} else {
		// Если запрос пустой, получаем все документы компании/папки
		if folderID != nil {
			documents, err = s.documentRepo.ListByFolder(ctx, *folderID)
		} else {
			documents, err = s.documentRepo.ListByCompany(ctx, companyID)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get documents: %w", err)
		}
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

// searchInElasticsearch выполняет полнотекстовой поиск в Elasticsearch
func (s *documentService) searchInElasticsearch(ctx context.Context, companyID uuid.UUID, query string, folderID *uuid.UUID) ([]uuid.UUID, error) {
	// Строим фильтры (must) - обязательные условия
	must := []map[string]interface{}{
		{
			"term": map[string]interface{}{
				"company_id": companyID.String(),
			},
		},
	}

	// Добавляем фильтр по папке, если указана
	if folderID != nil {
		must = append(must, map[string]interface{}{
			"term": map[string]interface{}{
				"folder_id": folderID.String(),
			},
		})
	}

	// Полнотекстовый поиск по имени и содержимому документа
	should := []map[string]interface{}{
		{
			"match": map[string]interface{}{
				"name": map[string]interface{}{
					"query": query,
					"boost": 2.0, // Повышаем релевантность совпадений в имени
				},
			},
		},
		{
			"match": map[string]interface{}{
				"text": map[string]interface{}{
					"query": query,
				},
			},
		},
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":                 must,
				"should":               should,
				"minimum_should_match": 1,
			},
		},
		"size": 1000, // Максимальное количество результатов
	}

	// Сериализуем запрос в JSON
	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Логируем запрос для отладки
	fmt.Printf("Elasticsearch query: %s\n", string(queryJSON))

	// Выполняем поиск
	res, err := s.elasticsearchClient.Search(
		s.elasticsearchClient.Search.WithContext(ctx),
		s.elasticsearchClient.Search.WithIndex("documents"),
		s.elasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
		s.elasticsearchClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("elasticsearch search error: %s", string(bodyBytes))
	}

	// Парсим результаты
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse elasticsearch response: %w", err)
	}

	// Извлекаем ID документов из результатов
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid elasticsearch response format")
	}

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid elasticsearch hits format")
	}

	fmt.Printf("Elasticsearch found %d documents\n", len(hitsArray))

	documentIDs := make([]uuid.UUID, 0, len(hitsArray))
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		docIDStr, ok := source["document_id"].(string)
		if !ok {
			continue
		}

		docID, err := uuid.Parse(docIDStr)
		if err != nil {
			continue
		}

		documentIDs = append(documentIDs, docID)
	}

	return documentIDs, nil
}

// getDocumentsByIDs получает документы по списку ID
func (s *documentService) getDocumentsByIDs(ctx context.Context, ids []uuid.UUID) ([]*ent.Document, error) {
	documents := make([]*ent.Document, 0, len(ids))
	for _, id := range ids {
		doc, err := s.documentRepo.GetByID(ctx, id)
		if err != nil {
			// Пропускаем документы, которые не удалось получить
			continue
		}
		documents = append(documents, doc)
	}
	return documents, nil
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

// ExtractAndIndexText извлекает текст из документа и индексирует его в Elasticsearch
// Использует docconv для извлечения текста из различных форматов документов
// Сохраняет извлеченный текст в индекс "documents" в Elasticsearch
func (s *documentService) ExtractAndIndexText(ctx context.Context, documentID uuid.UUID) error {
	// Получаем документ из БД
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Проверяем что документ поддерживает извлечение текста
	if !s.isExtractableText(document.MimeType) {
		return fmt.Errorf("document type %s does not support text extraction", document.MimeType)
	}

	// Скачиваем оригинальный файл из MinIO
	object, err := s.minioClient.GetObject(ctx, s.bucketName, document.FilePath, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to get file from minio: %w", err)
	}
	defer object.Close()

	// Извлекаем текст с помощью docconv
	c := client.New()
	convRes, err := c.Convert(object, document.Name)
	if err != nil {
		return fmt.Errorf("failed to extract text from document: %w", err)
	}

	// Очищаем текст от лишних пробелов и переносов строк
	extractedText := strings.TrimSpace(convRes.Body)
	if extractedText == "" {
		return fmt.Errorf("no text extracted from document")
	}

	// Подготавливаем документ для индексации в Elasticsearch
	// Индекс: documents
	// ID документа: UUID документа
	esDocument := map[string]interface{}{
		"document_id": document.ID.String(),
		"company_id":  document.CompanyID.String(),
		"folder_id":   nil,
		"name":        document.Name,
		"text":        extractedText,
		"mime_type":   document.MimeType,
		"file_size":   document.FileSize,
		"indexed_at":  time.Now().Format(time.RFC3339),
	}

	if document.FolderID != nil {
		esDocument["folder_id"] = document.FolderID.String()
	}

	// Сериализуем в JSON
	docJSON, err := json.Marshal(esDocument)
	if err != nil {
		return fmt.Errorf("failed to marshal document for elasticsearch: %w", err)
	}
	// Индексируем документ в Elasticsearch
	req := esapi.IndexRequest{
		Index:      "documents",
		DocumentID: document.ID.String(),
		Body:       bytes.NewReader(docJSON),
		Refresh:    "true",
	}

	esRes, err := req.Do(ctx, s.elasticsearchClient)
	if err != nil {
		return fmt.Errorf("failed to index document in elasticsearch: %w", err)
	}
	defer esRes.Body.Close()

	if esRes.IsError() {
		return fmt.Errorf("elasticsearch indexing error: %s", esRes.String())
	}

	fmt.Printf("Successfully indexed document %s in Elasticsearch index 'documents'\n", document.ID)
	return nil
}

// isExtractableText проверяет, можно ли извлечь текст из документа
func (s *documentService) isExtractableText(mimeType string) bool {
	extractableTypes := []string{
		// PDF
		"application/pdf",

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

		// Plain text
		"text/plain",
		"text/csv",

		// HTML
		"text/html",

		// Images (если docconv поддерживает OCR)
		"image/jpeg",
		"image/png",
		"image/tiff",
	}

	for _, et := range extractableTypes {
		if strings.EqualFold(mimeType, et) {
			return true
		}
	}

	return false
}
