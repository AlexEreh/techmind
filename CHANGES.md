# Изменения от 28 ноября 2024

## 1. Исправлена проблема с отображением документов

### Проблема
Документы не отображались в интерфейсе, хотя присутствовали в базе данных. Причина: API не передавал поле `name` в JSON-ответах.

### Решение

#### Backend (Go API)

1. **Добавлен общий тип ErrorResponse** (`/apps/api/internal/transport/http/handlers/dto.go`)
   - Создан общий тип `ErrorResponse` для унификации ответов с ошибками
   - Исправлены ошибки компиляции во всех handlers

2. **Обновлена структура DocumentResponse** (`/apps/api/internal/transport/http/handlers/document/dto.go`)
   - Добавлено поле `Name string` с JSON-тегом `json:"name"`
   - Поле размещено после `SenderID` и до `FilePath` в соответствии с логикой структуры

3. **Обновлены все document handlers** для передачи имени документа:
   - `get_by_company.go` - теперь возвращает `Name: docWithTags.Document.Name`
   - `get_by_folder.go` - теперь возвращает `Name: docWithTags.Document.Name`
   - `get_by_id.go` - теперь возвращает `Name: docWithTags.Document.Name`
   - `search.go` - теперь возвращает `Name: docWithTags.Document.Name`
   - `upload.go` - теперь возвращает `Name: document.Name`
   - `update.go` - теперь возвращает `Name: document.Name`

#### Frontend (TypeScript/React)

Фронтенд уже был готов к получению поля `name`:
- Интерфейс `Document` в `/apps/frontend/lib/api/types.ts` содержал поле `name: string`
- Компоненты корректно используют `document.name` для отображения

### Результат

✅ Документы теперь корректно отображаются в интерфейсе с их именами  
✅ API возвращает полную информацию о документах, включая имена файлов  
✅ Исправлены ошибки компиляции, связанные с отсутствующим `ErrorResponse`  
✅ К��д успешно компилируется без ошибок

### Тестирование

Для проверки работы:
1. Запустите API: `cd apps/api && go run cmd/main.go`
2. Запустите фронтенд: `cd apps/frontend && bun dev`
3. Откройте страницу `/files`
4. Выберите папку или корневую директорию
5. Документы должны отображаться с их именами

### API Endpoints, которые теперь возвращают поле `name`:

- `GET /private/documents/company/{company_id}` - все документы компании
- `GET /private/documents/folder/{folder_id}` - документы в папке
- `GET /private/documents/{id}` - конкретный документ
- `POST /private/documents/search` - поиск документов
- `POST /private/documents` - загрузка документа (возвращает созданный документ)
- `PUT /private/documents/{id}` - обновление документа

### Пример JSON-ответа:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "550e8400-e29b-41d4-a716-446655440001",
  "folder_id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "my-document.pdf",
  "file_path": "documents/550e8400-e29b-41d4-a716-446655440000.pdf",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "checksum": "abc123def456",
  "created_at": "2024-11-28T15:04:05Z",
  "tags": []
}
```

---

## 2. Исправлен предпросмотр документов

### Проблемы
1. Ошибка в консоли: "Unknown event handler property onPress"
2. API не возвращал `download_url`, из-за чего предпросмотр и скачивание не работали

### Решение

#### Backend (Go API)

1. **Обновлена структура DocumentWithTags** (`/apps/api/internal/service/service.go`)
   - Добавлено поле `DownloadURL string` для хранения ссылки на скачивание

2. **Обновлен сервисный слой** (`/apps/api/internal/service/document/document.go`)
   - `GetByID()` - добавлена генерация download URL
   - `GetByFolder()` - добавлена генерация download URL для каждого документа
   - `GetByCompany()` - добавлена генерация download URL для каждого документа
   - `Search()` - добавлена генерация download URL в результатах поиска

3. **Обновлены HTTP handlers** для передачи `DownloadURL`:
   - `get_by_id.go` - добавлено поле `DownloadURL`
   - `get_by_company.go` - добавлено поле `DownloadURL`
   - `get_by_folder.go` - добавлено поле `DownloadURL`
   - `search.go` - добавлено поле `DownloadURL`

#### Frontend (TypeScript/React)

1. **FileInfo.tsx**
   - Исправлено: `onPress` → `onClick` для компонента Chip (HeroUI не поддерживает onPress)
   - Кнопка "Скачать файл" теперь использует `document.download_url`

2. **FilePreview.tsx**
   - Предпросмотр изображений использует `download_url` если нет `preview_url`
   - Предпросмотр PDF использует `download_url` в iframe

### Результат

✅ Предпросмотр изображений работает  
✅ Предпросмотр PDF работает  
✅ Кнопка скачивания работает для всех типов файлов  
✅ Исправлена ошибка "Unknown event handler property onPress"  
✅ API возвращает все необходимые URL (preview_url и download_url)

### Как работает предпросмотр:

**Изображения:**
- Если есть `preview_url` - отображается превью
- Если нет `preview_url`, но есть `download_url` - отображается оригинал

**PDF:**
- Отображается в iframe используя `download_url`

**Остальные файлы:**
- Доступна кнопка "Скачать файл" с `download_url`

**Presigned URLs:**
- Генерируются через MinIO
- Срок действия: 1 час
- Безопасный доступ без авторизации

---

## 3. Улучшен UI дерева папок

### Проблемы
1. Элемент "Корень" имел светлый фон (bg-blue-100) в темной теме
2. У элемента "Корень" не было отступа слева, иконка папки была не выровнена

### Решение

#### Frontend (TypeScript/React)

**FolderTree.tsx:**
1. Заменен `bg-blue-100` на `bg-primary/20` для темной темы
   - Применено к элементу "Корень"
   - Применено ко всем папкам при выборе

2. Добавлен hover-эффект `hover:bg-default-100` для всех элементов

3. Добавлен невидимый спейсер для выравни��ания "Корня"
   - `<span className="w-6" />` компен��ирует отсутствие иконки шеврона
   - Иконка папки теперь выровнена с остальными папками

### Результат

✅ Единый стиль для всех элементов дер��ва папок  
✅ Правильные цвета для темной темы  
✅ Корректное выравнивание иконок  
✅ Улучшенная визуальная обратная связь при hover

---

## Пример обновленного JSON-ответа:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "company_id": "550e8400-e29b-41d4-a716-446655440001",
  "folder_id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "my-document.pdf",
  "file_path": "documents/550e8400-e29b-41d4-a716-446655440000.pdf",
  "preview_url": "https://minio.example.com/bucket/preview.jpg?X-Amz-...",
  "download_url": "https://minio.example.com/bucket/document.pdf?X-Amz-...",
  "file_size": 1024000,
  "mime_type": "application/pdf",
  "checksum": "abc123def456",
  "created_at": "2024-11-28T15:04:05Z",
  "tags": []
}
```


---
## 4. Исправлен поиск по тегам
### Проблема
На странице поиска клики по тегам не работали.
### Решение
Исправлен onClick для тегов в SearchPanel.tsx
### Результат
✅ Поиск по тегам работает
✅ Визуальная индикация выбранных тегов
✅ AND логика для множественного выбора

