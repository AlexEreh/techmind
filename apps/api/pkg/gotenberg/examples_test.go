package gotenberg_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"techmind/pkg/gotenberg"
)

// Этот файл содержит примеры использования клиента Gotenberg

func TestExample_ConvertURLToPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Создание клиента
	client := gotenberg.NewClient("http://localhost:3000")

	// Конвертация URL в PDF
	resp, err := client.ConvertURLToPDF(context.Background(), "https://example.com", &gotenberg.ChromiumRequest{
		PaperWidth:      "8.5",
		PaperHeight:     "11",
		MarginTop:       "0.5",
		MarginBottom:    "0.5",
		MarginLeft:      "0.5",
		MarginRight:     "0.5",
		PrintBackground: true,
		Landscape:       false,
		Scale:           1.0,
		OutputFilename:  "example",
	})

	if err != nil {
		t.Fatalf("failed to convert URL to PDF: %v", err)
	}

	// Сохранение результата
	if err := os.WriteFile("example.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("PDF created successfully: %s", resp.Filename)
}

func TestExample_ConvertHTMLToPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	htmlContent := []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test PDF</title>
			<style>
				body { font-family: Arial, sans-serif; }
				h1 { color: #333; }
			</style>
		</head>
		<body>
			<h1>Hello from Gotenberg!</h1>
			<p>This is a test PDF generated from HTML.</p>
			<img src="logo.png" alt="Logo" />
		</body>
		</html>
	`)

	// Загружаем изображение
	logoContent, err := os.ReadFile("testdata/logo.png")
	if err != nil {
		t.Skip("logo.png not found, skipping test")
	}

	resp, err := client.ConvertHTMLToPDF(
		context.Background(),
		gotenberg.File{Name: "index.html", Content: htmlContent},
		[]gotenberg.File{
			{Name: "logo.png", Content: logoContent},
		},
		&gotenberg.ChromiumRequest{
			PrintBackground: true,
			OutputFilename:  "test",
		},
	)

	if err != nil {
		t.Fatalf("failed to convert HTML to PDF: %v", err)
	}

	if err := os.WriteFile("test.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("PDF created successfully")
}

func TestExample_ConvertMarkdownToPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	htmlWrapper := []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Markdown PDF</title>
			<style>
				body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; }
			</style>
		</head>
		<body>
			{{ toHTML "content.md" }}
		</body>
		</html>
	`)

	markdownContent := []byte(`
# Hello from Markdown!

This is a **test** PDF generated from *Markdown*.

## Features

- Easy to write
- Easy to read
- Supports **formatting**

### Code Example

` + "```go" + `
func main() {
	fmt.Println("Hello, World!")
}
` + "```" + `
	`)

	resp, err := client.ConvertMarkdownToPDF(
		context.Background(),
		gotenberg.File{Name: "index.html", Content: htmlWrapper},
		[]gotenberg.File{
			{Name: "content.md", Content: markdownContent},
		},
		nil,
		&gotenberg.ChromiumRequest{
			OutputFilename: "markdown",
		},
	)

	if err != nil {
		t.Fatalf("failed to convert Markdown to PDF: %v", err)
	}

	if err := os.WriteFile("markdown.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("PDF created successfully")
}

func TestExample_ConvertOfficeToPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	// Загружаем Word документ
	docxContent, err := os.ReadFile("testdata/document.docx")
	if err != nil {
		t.Skip("document.docx not found, skipping test")
	}

	resp, err := client.ConvertOfficeToPDF(
		context.Background(),
		[]gotenberg.File{
			{Name: "document.docx", Content: docxContent},
		},
		&gotenberg.LibreOfficeRequest{
			Landscape:        false,
			NativePageRanges: "1-5",
			OutputFilename:   "office-doc",
		},
	)

	if err != nil {
		t.Fatalf("failed to convert Office document to PDF: %v", err)
	}

	if err := os.WriteFile("office-doc.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("PDF created successfully")
}

func TestExample_MergePDFs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	// Загружаем несколько PDF файлов
	pdf1, err := os.ReadFile("testdata/doc1.pdf")
	if err != nil {
		t.Skip("doc1.pdf not found, skipping test")
	}

	pdf2, err := os.ReadFile("testdata/doc2.pdf")
	if err != nil {
		t.Skip("doc2.pdf not found, skipping test")
	}

	resp, err := client.MergePDFs(
		context.Background(),
		[]gotenberg.File{
			{Name: "doc1.pdf", Content: pdf1},
			{Name: "doc2.pdf", Content: pdf2},
		},
		&gotenberg.PDFEnginesRequest{
			OutputFilename: "merged",
		},
	)

	if err != nil {
		t.Fatalf("failed to merge PDFs: %v", err)
	}

	if err := os.WriteFile("merged.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("PDFs merged successfully")
}

func TestExample_SplitPDF(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	pdfContent, err := os.ReadFile("testdata/multi-page.pdf")
	if err != nil {
		t.Skip("multi-page.pdf not found, skipping test")
	}

	// Разделение на интервалы
	resp, err := client.SplitPDFs(
		context.Background(),
		[]gotenberg.File{
			{Name: "multi-page.pdf", Content: pdfContent},
		},
		&gotenberg.SplitRequest{
			SplitMode: "intervals",
			SplitSpan: "2",
		},
	)

	if err != nil {
		t.Fatalf("failed to split PDF: %v", err)
	}

	if err := os.WriteFile("split-result.zip", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write result: %v", err)
	}

	t.Logf("PDF split successfully")
}

func TestExample_PDFMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	pdfContent, err := os.ReadFile("testdata/sample.pdf")
	if err != nil {
		t.Skip("sample.pdf not found, skipping test")
	}

	// Чтение метаданных
	metadata, err := client.ReadPDFMetadata(
		context.Background(),
		[]gotenberg.File{
			{Name: "sample.pdf", Content: pdfContent},
		},
	)

	if err != nil {
		t.Fatalf("failed to read metadata: %v", err)
	}

	t.Logf("Metadata: %+v", metadata)

	// Запись метаданных
	newMetadata := map[string]interface{}{
		"Author":   "Test Author",
		"Title":    "Test Document",
		"Subject":  "Testing",
		"Keywords": []string{"test", "example"},
	}

	resp, err := client.WritePDFMetadata(
		context.Background(),
		[]gotenberg.File{
			{Name: "sample.pdf", Content: pdfContent},
		},
		newMetadata,
		"with-metadata",
	)

	if err != nil {
		t.Fatalf("failed to write metadata: %v", err)
	}

	if err := os.WriteFile("with-metadata.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("Metadata written successfully")
}

func TestExample_ScreenshotURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	resp, err := client.ScreenshotURL(
		context.Background(),
		"https://example.com",
		&gotenberg.ScreenshotRequest{
			Width:            1920,
			Height:           1080,
			Format:           "png",
			OptimizeForSpeed: true,
			OutputFilename:   "screenshot",
		},
	)

	if err != nil {
		t.Fatalf("failed to take screenshot: %v", err)
	}

	if err := os.WriteFile("screenshot.png", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write screenshot: %v", err)
	}

	t.Logf("Screenshot created successfully")
}

func TestExample_WithOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Создание клиента с опциями
	client := gotenberg.NewClient(
		"http://localhost:3000",
		gotenberg.WithTrace("my-trace-id"),
		gotenberg.WithHTTPClient(&http.Client{
			Timeout: 60 * time.Second,
		}),
	)

	resp, err := client.GetVersion(context.Background())
	if err != nil {
		t.Fatalf("failed to get version: %v", err)
	}

	t.Logf("Gotenberg version: %s", resp)
}

func TestExample_HealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	health, err := client.GetHealth(context.Background())
	if err != nil {
		t.Fatalf("failed to check health: %v", err)
	}

	t.Logf("Health status: %s", health.Status)
	for name, check := range health.Details {
		t.Logf("  %s: %s (checked at %v)", name, check.Status, check.Timestamp)
	}
}

func TestExample_AdvancedChromium(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := gotenberg.NewClient("http://localhost:3000")

	// Продвинутая конфигурация с заголовками, footer, cookies и т.д.
	headerHTML := []byte(`
		<html>
		<head>
			<style>
				body { font-size: 10px; text-align: center; }
			</style>
		</head>
		<body>
			<p>Page <span class="pageNumber"></span> of <span class="totalPages"></span></p>
		</body>
		</html>
	`)

	footerHTML := []byte(`
		<html>
		<head>
			<style>
				body { font-size: 10px; text-align: center; }
			</style>
		</head>
		<body>
			<p>Generated on <span class="date"></span></p>
		</body>
		</html>
	`)

	resp, err := client.ConvertURLToPDF(
		context.Background(),
		"https://example.com",
		&gotenberg.ChromiumRequest{
			// Page properties
			PaperWidth:              "8.27", // A4
			PaperHeight:             "11.7", // A4
			MarginTop:               "1",
			MarginBottom:            "1",
			MarginLeft:              "1",
			MarginRight:             "1",
			PrintBackground:         true,
			Landscape:               false,
			Scale:                   1.0,
			GenerateDocumentOutline: true,
			GenerateTaggedPDF:       true,

			// Header & Footer
			Header: &gotenberg.File{Name: "header.html", Content: headerHTML},
			Footer: &gotenberg.File{Name: "footer.html", Content: footerHTML},

			// Wait before rendering
			WaitDelay: "2s",

			// Cookies
			Cookies: []gotenberg.Cookie{
				{
					Name:   "session_id",
					Value:  "abc123",
					Domain: "example.com",
					Secure: true,
				},
			},

			// Custom headers
			UserAgent: "Mozilla/5.0 (Custom)",
			ExtraHTTPHeaders: map[string]string{
				"Authorization": "Bearer token123",
			},

			// Error handling
			FailOnHTTPStatusCodes:       []int{499},
			FailOnConsoleExceptions:     true,
			FailOnResourceLoadingFailed: false,

			// PDF/A
			PDFA:  "PDF/A-1b",
			PDFUA: true,

			// Metadata
			Metadata: map[string]interface{}{
				"Author":  "John Doe",
				"Title":   "Advanced PDF",
				"Subject": "Testing advanced features",
			},

			// Security
			UserPassword:  "user123",
			OwnerPassword: "owner456",

			OutputFilename: "advanced",
		},
	)

	if err != nil {
		t.Fatalf("failed to convert: %v", err)
	}

	if err := os.WriteFile("advanced.pdf", resp.Body, 0644); err != nil {
		t.Fatalf("failed to write PDF: %v", err)
	}

	t.Logf("Advanced PDF created successfully")
}
