# Интеграция Gotenberg и генерация PDF превью

## Сводка изменений

Была реализована полная интеграция Gotenberg для конвертации документов в PDF превью.

## Файлы изменены

### 1. Конфигурация

#### `pkg/config/config.go`
- ✅ Добавлена секция `Gotenberg` с полями:
  - `URL` - адрес Gotenberg API
  - `Enabled` - флаг включения/выключения
  - `Timeout` - таймаут запросов в секундах

#### `dev.yml`
- ✅ Добавлена конфигурация Gotenberg:
```yaml
gotenberg:
  url: "http://localhost:3000"
  enabled: true
  timeout: 60
```

#### `docker-compose.dev.yml`
- ✅ Исправлен сервис Gotenberg с правильным маппингом портов
- ✅ Добавлены команды для настройки таймаута и логирования

### 2. Подключение

#### `internal/connections/gotenberg/gotenberg.go` (НОВЫЙ)
- ✅ Создан провайдер для Uber FX
- ✅ Инициализация Gotenberg клиента с настройками из конфига
- ✅ Поддержка отключения через флаг `enabled`

#### `app/app.go`
- ✅ Добавлен импорт `gotenberg` connection
- ✅ Зарегистрирован `gotenberg.New` в `fx.Provide`

### 3. Сервис документов

#### `internal/service/document/document.go`
- ✅ Добавлены импорты `bytes`, `strings`, `gotenberg`
- ✅ Добавлены поля в структуру `documentService`:
  - `gotenbergClient *gotenberg.Client`
  - `gotenbergEnabled bool`
- ✅ Обновлен конструктор `NewService` с параметром `gotenbergClient`
- ✅ Реализована функция `GeneratePDFPreview()` - основная функция конвертации
- ✅ Добавлена функция `isConvertibleToPDF()` - проверка поддерживаемых типов
- ✅ Добавлена функция `isOfficeDocument()` - проверка Office документов

#### `internal/service/service.go`
- ✅ Добавлен метод `GeneratePDFPreview` в интерфейс `DocumentService`

### 4. Документация

#### `docs/PREVIEW_GENERATION.md` (НОВЫЙ)
- ✅ Полное описание функции генерации превью
- ✅ Поддерживаемые форматы файлов
- ✅ Примеры использования
- ✅ Архитектура и процесс генерации
- ✅ Обработка ошибок
- ✅ Настройка и troubleshooting

#### `docs/GOTENBERG_SETUP.md` (НОВЫЙ)
- ✅ Инструкции по настройке Gotenberg
- ✅ Конфигурация для dev/prod
- ✅ Docker Compose настройки
- ✅ Troubleshooting

#### `internal/service/document/example_preview_generation.go` (НОВЫЙ)
- ✅ Примеры использования функции
- ✅ Пример массовой генерации превью
- ✅ Пример проверки поддерживаемых типов

## Функциональность

### GeneratePDFPreview

Основная функция для конвертации документов:

```go
func (s *documentService) GeneratePDFPreview(ctx context.Context, documentID uuid.UUID) error
```

**Что делает:**
1. Проверяет доступность Gotenberg
2. Загружает документ из БД
3. Проверяет поддержку формата
4. Скачивает оригинальный файл из MinIO
5. Конвертирует через Gotenberg API
6. Загружает PDF в MinIO (`{company_id}/previews/{preview_id}.pdf`)
7. Обновляет `preview_file_path` в БД
8. При ошибке откатывает изменения

**Поддерживаемые форматы:**
- Microsoft Office: .docx, .xlsx, .pptx, .doc, .xls, .ppt
- OpenDocument: .odt, .ods, .odp
- Другие: .rtf, .html

## Использование

### Простой пример

```go
// Генерация превью для документа
err := documentService.GeneratePDFPreview(ctx, documentID)
if err != nil {
    log.Printf("Error: %v", err)
}

// Получение URL превью
previewURL, err := documentService.GetPreviewURL(ctx, documentID)
```

### Асинхронная генерация

```go
// После загрузки документа
document, err := docService.Upload(ctx, input)
if err != nil {
    return nil, err
}

// Генерируем превью в фоне
go func() {
    err := docService.GeneratePDFPreview(context.Background(), document.ID)
    if err != nil {
        log.Printf("Preview generation failed: %v", err)
    }
}()
```

## Запуск

### Локальная разработка

```bash
# 1. Запустить сервисы
docker-compose -f docker-compose.dev.yml up -d

# 2. Проверить Gotenberg
curl http://localhost:3000/health

# 3. Запустить приложение
cd apps/api
go run cmd/main.go
```

### Проверка работы

```bash
# Health check
curl http://localhost:3000/health

# Version check
curl http://localhost:3000/version
```

## Отключение Gotenberg

Если Gotenberg не нужен, можно отключить в конфиге:

```yaml
gotenberg:
  url: "http://localhost:3000"
  enabled: false  # Отключить
  timeout: 60
```

При отключении:
- Клиент не создается (`nil`)
- Функция `GeneratePDFPreview` вернет ошибку
- Остальные функции работают нормально

## Тестирование

```bash
# Тесты сервиса документов
cd apps/api
go test ./internal/service/document/... -v

# Проверка компиляции
go build ./...
```

## Следующие шаги

1. ✅ Добавить HTTP endpoint для генерации превью
2. ✅ Реализовать фоновую обработку через очередь
3. ✅ Добавить метрики и мониторинг
4. ✅ Настроить автоматическую генерацию при загрузке

## Полезные ссылки

- [Gotenberg Documentation](https://gotenberg.dev/)
- [MinIO Go Client](https://min.io/docs/minio/linux/developers/go/minio-go.html)
- [Uber FX](https://uber-go.github.io/fx/)

