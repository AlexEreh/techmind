# Gotenberg Go Client

Полнофункциональная библиотека на Go для работы с [Gotenberg API](https://gotenberg.dev) - сервисом для конвертации документов в PDF и создания скриншотов.

## Установка

```bash
# Библиотека является частью проекта techmind
import "techmind/pkg/gotenberg"
```

## Быстрый старт

```go
package main

import (
    "context"
    "os"
    
    "techmind/pkg/gotenberg"
)

func main() {
    // Создание клиента
    client := gotenberg.NewClient("http://localhost:3000")
    
    // Конвертация URL в PDF
    resp, err := client.ConvertURLToPDF(
        context.Background(),
        "https://example.com",
        &gotenberg.ChromiumRequest{
            PrintBackground: true,
            OutputFilename:  "example",
        },
    )
    if err != nil {
        panic(err)
    }
    
    // Сохранение результата
    os.WriteFile("example.pdf", resp.Body, 0644)
}
```

## Возможности

### Конвертация с помощью Chromium

#### Конвертация URL в PDF

```go
resp, err := client.ConvertURLToPDF(ctx, "https://example.com", &gotenberg.ChromiumRequest{
    PaperWidth:      "8.5",
    PaperHeight:     "11",
    MarginTop:       "0.5",
    MarginBottom:    "0.5",
    MarginLeft:      "0.5",
    MarginRight:     "0.5",
    PrintBackground: true,
    Landscape:       false,
    Scale:           1.0,
})
```

#### Конвертация HTML в PDF

```go
htmlContent := []byte(`<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><h1>Hello!</h1></body>
</html>`)

resp, err := client.ConvertHTMLToPDF(
    ctx,
    gotenberg.File{Name: "index.html", Content: htmlContent},
    nil, // дополнительные файлы (CSS, изображения и т.д.)
    &gotenberg.ChromiumRequest{
        PrintBackground: true,
    },
)
```

#### Конвертация Markdown в PDF

```go
htmlWrapper := []byte(`<!DOCTYPE html>
<html>
<body>{{ toHTML "content.md" }}</body>
</html>`)

markdownContent := []byte(`# Hello from Markdown!`)

resp, err := client.ConvertMarkdownToPDF(
    ctx,
    gotenberg.File{Name: "index.html", Content: htmlWrapper},
    []gotenberg.File{{Name: "content.md", Content: markdownContent}},
    nil, // дополнительные файлы
    nil, // опции
)
```

### Создание скриншотов

```go
resp, err := client.ScreenshotURL(ctx, "https://example.com", &gotenberg.ScreenshotRequest{
    Width:            1920,
    Height:           1080,
    Format:           "png", // "png", "jpeg", "webp"
    Quality:          100,   // для JPEG
    OptimizeForSpeed: true,
})
```

### Конвертация Office документов

```go
docxContent, _ := os.ReadFile("document.docx")

resp, err := client.ConvertOfficeToPDF(
    ctx,
    []gotenberg.File{{Name: "document.docx", Content: docxContent}},
    &gotenberg.LibreOfficeRequest{
        Landscape:        false,
        NativePageRanges: "1-5",
        Merge:            false,
    },
)
```

Поддерживаемые форматы:
- Word: `.doc`, `.docx`, `.dot`, `.dotx`, `.docm`, `.dotm`
- Excel: `.xls`, `.xlsx`, `.xlsm`, `.xlt`, `.xltx`, `.xltm`
- PowerPoint: `.ppt`, `.pptx`, `.pptm`, `.pot`, `.potx`, `.potm`
- И многие другие (более 100 форматов)

### Работа с PDF

#### Объединение PDF

```go
pdf1, _ := os.ReadFile("doc1.pdf")
pdf2, _ := os.ReadFile("doc2.pdf")

resp, err := client.MergePDFs(
    ctx,
    []gotenberg.File{
        {Name: "doc1.pdf", Content: pdf1},
        {Name: "doc2.pdf", Content: pdf2},
    },
    &gotenberg.PDFEnginesRequest{
        PDFA: "PDF/A-1b", // опционально
    },
)
```

#### Разделение PDF

```go
pdfContent, _ := os.ReadFile("multi-page.pdf")

// По интервалам
resp, err := client.SplitPDFs(
    ctx,
    []gotenberg.File{{Name: "multi-page.pdf", Content: pdfContent}},
    &gotenberg.SplitRequest{
        SplitMode: "intervals",
        SplitSpan: "2", // каждые 2 страницы
    },
)

// По диапазонам страниц
resp, err := client.SplitPDFs(
    ctx,
    []gotenberg.File{{Name: "multi-page.pdf", Content: pdfContent}},
    &gotenberg.SplitRequest{
        SplitMode:  "pages",
        SplitSpan:  "1-3,5,8-10",
        SplitUnify: true, // объединить в один файл
    },
)
```

#### Чтение и запись метаданных

```go
// Чтение метаданных
metadata, err := client.ReadPDFMetadata(ctx, []gotenberg.File{
    {Name: "sample.pdf", Content: pdfContent},
})

// Запись метаданных
newMetadata := map[string]interface{}{
    "Author":   "John Doe",
    "Title":    "My Document",
    "Subject":  "Important",
    "Keywords": []string{"pdf", "metadata"},
}

resp, err := client.WritePDFMetadata(
    ctx,
    []gotenberg.File{{Name: "sample.pdf", Content: pdfContent}},
    newMetadata,
    "output",
)
```

#### Конвертация в PDF/A и PDF/UA

```go
resp, err := client.ConvertToPDFA(
    ctx,
    []gotenberg.File{{Name: "document.pdf", Content: pdfContent}},
    "PDF/A-1b", // "PDF/A-1b", "PDF/A-2b", "PDF/A-3b"
    true,       // PDF/UA для доступности
    "output",
)
```

#### Сглаживание PDF

```go
resp, err := client.FlattenPDFs(
    ctx,
    []gotenberg.File{{Name: "with-forms.pdf", Content: pdfContent}},
    "flattened",
)
```

#### Шифрование PDF

```go
resp, err := client.EncryptPDFs(
    ctx,
    []gotenberg.File{{Name: "document.pdf", Content: pdfContent}},
    "user_password",  // пароль для открытия
    "owner_password", // пароль для полного доступа
    "encrypted",
)
```

#### Встраивание файлов в PDF

```go
xmlContent, _ := os.ReadFile("invoice.xml")

resp, err := client.EmbedFilesInPDF(
    ctx,
    []gotenberg.File{{Name: "invoice.pdf", Content: pdfContent}},
    []gotenberg.File{{Name: "invoice.xml", Content: xmlContent}},
    "with-attachments",
)
```

## Продвинутые возможности

### Настройка страницы

```go
req := &gotenberg.ChromiumRequest{
    // Размер бумаги
    PaperWidth:  "8.27", // A4 ширина в дюймах
    PaperHeight: "11.7", // A4 высота в дюймах
    
    // Поля (в дюймах, можно использовать: pt, px, in, mm, cm, pc)
    MarginTop:    "1in",
    MarginBottom: "1in",
    MarginLeft:   "1in",
    MarginRight:  "1in",
    
    // Ориентация
    Landscape: false,
    
    // Масштаб
    Scale: 1.0,
    
    // Диапазон страниц
    NativePageRanges: "1-5,8,11-13",
    
    // Одна страница
    SinglePage: false,
    
    // Предпочитать CSS размер
    PreferCSSPageSize: false,
    
    // Фон
    PrintBackground: true,
    OmitBackground:  false, // прозрачный фон
}
```

### Заголовки и подвалы

```go
headerHTML := []byte(`
<html>
<head>
    <style>
        body { font-size: 12px; text-align: center; }
    </style>
</head>
<body>
    <p>Страница <span class="pageNumber"></span> из <span class="totalPages"></span></p>
</body>
</html>
`)

req := &gotenberg.ChromiumRequest{
    Header: &gotenberg.File{Name: "header.html", Content: headerHTML},
    Footer: &gotenberg.File{Name: "footer.html", Content: footerHTML},
}
```

Доступные классы:
- `.pageNumber` - текущий номер страницы
- `.totalPages` - всего страниц
- `.date` - дата печати
- `.title` - заголовок документа
- `.url` - URL документа

### Ожидание перед рендерингом

```go
req := &gotenberg.ChromiumRequest{
    // Ждать 5 секунд
    WaitDelay: "5s",
    
    // Или ждать выполнения JavaScript выражения
    WaitForExpression: "window.status === 'ready'",
}
```

### Cookies и HTTP заголовки

```go
req := &gotenberg.ChromiumRequest{
    // Cookies
    Cookies: []gotenberg.Cookie{
        {
            Name:     "session_id",
            Value:    "abc123",
            Domain:   "example.com",
            Path:     "/",
            Secure:   true,
            HTTPOnly: true,
            SameSite: "Strict",
        },
    },
    
    // User Agent
    UserAgent: "Mozilla/5.0 (Custom Agent)",
    
    // Дополнительные заголовки
    ExtraHTTPHeaders: map[string]string{
        "Authorization": "Bearer token",
        "X-Custom":      "value",
    },
}
```

### Обработка ошибок

```go
req := &gotenberg.ChromiumRequest{
    // Ошибка при определенных HTTP статусах
    FailOnHTTPStatusCodes: []int{499, 599}, // 400-499, 500-599
    FailOnResourceHTTPStatusCodes: []int{499},
    
    // Ошибка при неудачной загрузке ресурсов
    FailOnResourceLoadingFailed: true,
    
    // Ошибка при исключениях в консоли
    FailOnConsoleExceptions: true,
}
```

### Режим производительности

```go
req := &gotenberg.ChromiumRequest{
    // Не ждать события "network idle" для ускорения
    SkipNetworkIdleEvent: true, // по умолчанию true
}
```

### Загрузка файлов из URL

```go
req := &gotenberg.ChromiumRequest{
    DownloadFrom: []gotenberg.DownloadFrom{
        {
            URL: "https://example.com/document.pdf",
            ExtraHTTPHeaders: map[string]string{
                "Authorization": "Bearer token",
            },
            Embedded: false, // встроить в результирующий PDF
        },
    },
}
```

## Конфигурация клиента

### Опции клиента

```go
import "net/http"
import "time"

client := gotenberg.NewClient(
    "http://localhost:3000",
    
    // Пользовательский HTTP клиент
    gotenberg.WithHTTPClient(&http.Client{
        Timeout: 60 * time.Second,
    }),
    
    // Trace ID для всех запросов
    gotenberg.WithTrace("my-trace-id"),
)
```

### Health Check

```go
health, err := client.GetHealth(ctx)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", health.Status)
for name, check := range health.Details {
    fmt.Printf("  %s: %s\n", name, check.Status)
}
```

### Версия

```go
version, err := client.GetVersion(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Gotenberg version: %s\n", version)
```

## Структура Response

```go
type Response struct {
    Body        []byte  // Содержимое файла
    ContentType string  // MIME тип
    Filename    string  // Имя файла из заголовка
    Trace       string  // Trace ID запроса
}
```

## Форматы PDF/A

- `PDF/A-1b` - базовый уровень PDF/A-1
- `PDF/A-2b` - базовый уровень PDF/A-2
- `PDF/A-3b` - базовый уровень PDF/A-3 (с поддержкой встраивания файлов)

## Обработка ошибок

Все методы возвращают ошибку при:
- Проблемах с сетью
- HTTP статусах >= 400
- Невалидных параметрах

```go
resp, err := client.ConvertURLToPDF(ctx, url, req)
if err != nil {
    // Обработка ошибки
    log.Printf("Failed to convert: %v", err)
    return
}

// Использование результата
os.WriteFile("output.pdf", resp.Body, 0644)
```

## Примеры

Дополнительные примеры использования смотрите в файле [examples_test.go](examples_test.go).

## Поддерживаемые форматы Office документов

LibreOffice поддерживает конвертацию из более чем 100 форматов:

**Документы:**
- Word: `.doc`, `.docx`, `.dot`, `.dotx`, `.docm`, `.dotm`, `.rtf`
- OpenDocument: `.odt`, `.ott`, `.fodt`
- Другие: `.wpd`, `.wps`, `.abw`, `.zabw`

**Таблицы:**
- Excel: `.xls`, `.xlsx`, `.xlsm`, `.xlt`, `.xltx`, `.xltm`
- OpenDocument: `.ods`, `.ots`, `.fods`
- Другие: `.csv`, `.dbf`, `.slk`

**Презентации:**
- PowerPoint: `.ppt`, `.pptx`, `.pptm`, `.pot`, `.potx`, `.potm`
- OpenDocument: `.odp`, `.otp`, `.fodp`
- Другие: `.key`, `.pps`

**Изображения:**
- `.bmp`, `.gif`, `.jpeg`, `.jpg`, `.png`, `.svg`, `.tiff`, `.webp`

И многие другие форматы.

## Лицензия

MIT

## Ссылки

- [Официальная документация Gotenberg](https://gotenberg.dev)
- [Gotenberg на GitHub](https://github.com/gotenberg/gotenberg)

