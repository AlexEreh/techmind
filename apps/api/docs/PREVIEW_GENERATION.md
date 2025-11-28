# Генерация PDF превью для документов

## Описание

Функция `GeneratePDFPreview` автоматически конвертирует загруженные документы в PDF формат и сохраняет их как превью в MinIO (S3). Это позволяет пользователям просматривать документы в браузере без необходимости скачивания оригинальных файлов.

## Возможности

- ✅ Конвертация Office документов (Word, Excel, PowerPoint)
- ✅ Поддержка старых форматов Office (.doc, .xls, .ppt)
- ✅ Конвертация OpenDocument форматов (.odt, .ods, .odp)
- ✅ Поддержка RTF и HTML файлов
- ✅ Автоматическое сохранение в MinIO
- ✅ Обновление ссылки на превью в базе данных
- ✅ Откат при ошибках (удаление загруженного файла, если не удалось обновить БД)

## Поддерживаемые форматы

### Microsoft Office
- `.docx` - Word документы (современный формат)
- `.xlsx` - Excel таблицы (современный формат)
- `.pptx` - PowerPoint презентации (современный формат)
- `.doc` - Word 97-2003
- `.xls` - Excel 97-2003
- `.ppt` - PowerPoint 97-2003

### OpenDocument
- `.odt` - OpenDocument текст
- `.ods` - OpenDocument таблицы
- `.odp` - OpenDocument презентации

### Другие
- `.rtf` - Rich Text Format
- `.html` - HTML документы

## Использование

### Базовое использование

```go
import (
    "context"
    "github.com/google/uuid"
)

// Генерация превью для конкретного документа
func generatePreview(docService service.DocumentService, documentID uuid.UUID) error {
    ctx := context.Background()
    
    err := docService.GeneratePDFPreview(ctx, documentID)
    if err != nil {
        return fmt.Errorf("failed to generate preview: %w", err)
    }
    
    // Получить URL превью
    previewURL, err := docService.GetPreviewURL(ctx, documentID)
    if err != nil {
        return fmt.Errorf("failed to get preview URL: %w", err)
    }
    
    fmt.Printf("Preview URL: %s\n", previewURL)
    return nil
}
```

### Интеграция в процесс загрузки

```go
// После загрузки документа автоматически генерируем превью
func uploadWithPreview(docService service.DocumentService, input service.DocumentUploadInput) (*ent.Document, error) {
    ctx := context.Background()
    
    // Загружаем документ
    document, err := docService.Upload(ctx, input)
    if err != nil {
        return nil, err
    }
    
    // Пытаемся сгенерировать превью (асинхронно или в фоне)
    go func() {
        err := docService.GeneratePDFPreview(context.Background(), document.ID)
        if err != nil {
            log.Printf("Failed to generate preview for document %s: %v", document.ID, err)
        }
    }()
    
    return document, nil
}
```

### Массовая генерация превью

```go
// Генерация превью для всех документов компании без превью
func batchGeneratePreviews(docService service.DocumentService, companyID uuid.UUID) error {
    ctx := context.Background()
    
    // Получаем все документы
    documents, err := docService.GetByCompany(ctx, companyID)
    if err != nil {
        return err
    }
    
    for _, doc := range documents {
        // Пропускаем документы, у которых уже есть превью
        if doc.Document.PreviewFilePath != nil {
            continue
        }
        
        // Генерируем превью
        if err := docService.GeneratePDFPreview(ctx, doc.Document.ID); err != nil {
            log.Printf("Error generating preview for %s: %v", doc.Document.Name, err)
            continue
        }
        
        log.Printf("Generated preview for: %s", doc.Document.Name)
    }
    
    return nil
}
```

## Архитектура

### Процесс генерации превью

1. **Проверка доступности Gotenberg**
   - Проверяет, что Gotenberg сервис настроен и доступен

2. **Получение документа**
   - Загружает метаданные документа из базы данных
   - Проверяет поддерживаемость формата файла

3. **Скачивание оригинального файла**
   - Загружает оригинальный файл из MinIO

4. **Конвертация в PDF**
   - Использует Gotenberg API для конвертации
   - LibreOffice используется для Office документов

5. **Загрузка превью**
   - Генерирует уникальное имя для PDF файла
   - Загружает PDF в MinIO в папку `{company_id}/previews/`

6. **Обновление базы данных**
   - Обновляет поле `preview_file_path` в таблице документов
   - При ошибке откатывает изменения (удаляет загруженный файл)

### Структура хранения в MinIO

```
documents/
├── {company_id}/
│   ├── {document_id}.docx         # Оригинальный файл
│   ├── {document_id}.xlsx
│   └── previews/
│       ├── {preview_id}.pdf       # PDF превью
│       └── {preview_id}.pdf
```

## Обработка ошибок

Функция обрабатывает следующие типы ошибок:

- ❌ **Gotenberg не настроен** - возвращает ошибку о недоступности сервиса
- ❌ **Документ не найден** - возвращает ошибку "document not found"
- ❌ **Неподдерживаемый формат** - возвращает ошибку с указанием MIME-типа
- ❌ **Ошибка скачивания из MinIO** - возвращает ошибку "failed to get file from minio"
- ❌ **Ошибка конвертации** - возвращает ошибку от Gotenberg
- ❌ **Ошибка загрузки превью** - возвращает ошибку "failed to upload preview to minio"
- ❌ **Ошибка обновления БД** - откатывает загрузку файла и возвращает ошибку

## Настройка Gotenberg

### Docker Compose

```yaml
services:
  gotenberg:
    image: gotenberg/gotenberg:7
    ports:
      - "3000:3000"
    environment:
      - DISABLE_GOOGLE_CHROME=1  # Опционально, если не нужен Chrome
    command:
      - "gotenberg"
      - "--api-timeout=30s"
      - "--log-level=info"
```

### Инициализация в приложении

```go
import (
    "net/http"
    "time"
    "techmind/pkg/gotenberg"
)

// Создание Gotenberg клиента
gotenbergClient := gotenberg.NewClient(
    "http://gotenberg:3000",
    gotenberg.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
)

// Инициализация сервиса документов с Gotenberg
documentService := document.NewService(
    documentRepo,
    documentTagRepo,
    tagRepo,
    folderRepo,
    minioClient,
    gotenbergClient, // передаем клиент
)
```

## Производительность

### Рекомендации

1. **Асинхронная генерация** - генерируйте превью в фоновом процессе, чтобы не блокировать загрузку документа
2. **Очередь задач** - используйте систему очередей (Redis, RabbitMQ) для обработки больших объемов
3. **Кэширование** - MinIO presigned URLs действуют 1 час, используйте их повторно
4. **Мониторинг** - отслеживайте время конвертации и размер файлов

### Типичное время конвертации

- **Word документ (10 страниц)**: ~2-3 секунды
- **Excel таблица (средняя)**: ~1-2 секунды
- **PowerPoint (20 слайдов)**: ~3-5 секунд

## API эндпоинты (пример)

```go
// POST /api/documents/{id}/preview
// Генерация превью для документа
func GeneratePreviewHandler(w http.ResponseWriter, r *http.Request) {
    documentID := getDocumentIDFromPath(r)
    
    err := documentService.GeneratePDFPreview(r.Context(), documentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    json.NewEncoder(w).Encode(map[string]string{
        "status": "success",
        "message": "Preview generated successfully",
    })
}
```

## Тестирование

```go
func TestGeneratePDFPreview(t *testing.T) {
    // Создаем тестовый документ
    doc := createTestDocument(t, "test.docx")
    
    // Генерируем превью
    err := docService.GeneratePDFPreview(context.Background(), doc.ID)
    assert.NoError(t, err)
    
    // Проверяем что preview path обновился
    updated, err := documentRepo.GetByID(context.Background(), doc.ID)
    assert.NoError(t, err)
    assert.NotNil(t, updated.PreviewFilePath)
    
    // Проверяем что файл существует в MinIO
    _, err = minioClient.StatObject(context.Background(), "documents", *updated.PreviewFilePath, minio.StatObjectOptions{})
    assert.NoError(t, err)
}
```

## Troubleshooting

### Gotenberg недоступен
```
Error: gotenberg is not enabled or configured
```
**Решение**: Убедитесь что Gotenberg запущен и передан в конструктор сервиса

### Неподдерживаемый формат
```
Error: document type application/octet-stream is not convertible to PDF
```
**Решение**: Проверьте что MIME-тип файла установлен правильно при загрузке

### Ошибка конвертации
```
Error: failed to convert office document to PDF
```
**Решение**: Проверьте логи Gotenberg, файл может быть поврежден

## См. также

- [Gotenberg Documentation](https://gotenberg.dev/)
- [MinIO Documentation](https://min.io/docs/)
- Пример использования: `example_preview_generation.go`

