package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"techmind/pkg/gotenberg"
)

// Пример простого использования Gotenberg клиента
func main() {
	// Создаем клиент с таймаутом 60 секунд
	client := gotenberg.NewClient(
		"http://localhost:3000",
		gotenberg.WithHTTPClient(&http.Client{
			Timeout: 60 * time.Second,
		}),
	)

	// Проверяем здоровье сервера
	if err := checkHealth(client); err != nil {
		log.Fatalf("Gotenberg server is not healthy: %v", err)
	}

	// Примеры использования
	if err := convertURL(client); err != nil {
		log.Printf("Failed to convert URL: %v", err)
	}

	if err := convertHTML(client); err != nil {
		log.Printf("Failed to convert HTML: %v", err)
	}

	if err := mergeDocuments(client); err != nil {
		log.Printf("Failed to merge documents: %v", err)
	}

	fmt.Println("All examples completed successfully!")
}

// checkHealth проверяет статус здоровья Gotenberg сервера
func checkHealth(client *gotenberg.Client) error {
	health, err := client.GetHealth(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Gotenberg Health Status: %s\n", health.Status)
	for name, check := range health.Details {
		fmt.Printf("  - %s: %s (checked at %v)\n", name, check.Status, check.Timestamp)
	}

	return nil
}

// convertURL конвертирует URL в PDF
func convertURL(client *gotenberg.Client) error {
	fmt.Println("\n=== Converting URL to PDF ===")

	resp, err := client.ConvertURLToPDF(
		context.Background(),
		"https://example.com",
		&gotenberg.ChromiumRequest{
			PaperWidth:      "8.5",
			PaperHeight:     "11",
			MarginTop:       "0.5",
			MarginBottom:    "0.5",
			MarginLeft:      "0.5",
			MarginRight:     "0.5",
			PrintBackground: true,
			OutputFilename:  "example-url",
		},
	)
	if err != nil {
		return err
	}

	filename := "example-url.pdf"
	if err := os.WriteFile(filename, resp.Body, 0644); err != nil {
		return err
	}

	fmt.Printf("✓ PDF created: %s (size: %d bytes)\n", filename, len(resp.Body))
	return nil
}

// convertHTML конвертирует HTML в PDF
func convertHTML(client *gotenberg.Client) error {
	fmt.Println("\n=== Converting HTML to PDF ===")

	htmlContent := []byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="utf-8">
			<title>Test Document</title>
			<style>
				body { 
					font-family: Arial, sans-serif; 
					padding: 20px;
				}
				h1 { 
					color: #2563eb; 
					border-bottom: 2px solid #2563eb;
					padding-bottom: 10px;
				}
				.info {
					background-color: #e0f2fe;
					border-left: 4px solid #0284c7;
					padding: 10px;
					margin: 20px 0;
				}
			</style>
		</head>
		<body>
			<h1>Hello from Gotenberg!</h1>
			<div class="info">
				<p><strong>Info:</strong> This PDF was generated from HTML using Gotenberg.</p>
			</div>
			<p>This is a test document with some <strong>bold text</strong> and <em>italic text</em>.</p>
			<ul>
				<li>Feature 1</li>
				<li>Feature 2</li>
				<li>Feature 3</li>
			</ul>
		</body>
		</html>
	`)

	resp, err := client.ConvertHTMLToPDF(
		context.Background(),
		gotenberg.File{Name: "index.html", Content: htmlContent},
		nil, // no additional files
		&gotenberg.ChromiumRequest{
			PrintBackground: true,
			OutputFilename:  "example-html",
		},
	)
	if err != nil {
		return err
	}

	filename := "example-html.pdf"
	if err := os.WriteFile(filename, resp.Body, 0644); err != nil {
		return err
	}

	fmt.Printf("✓ PDF created: %s (size: %d bytes)\n", filename, len(resp.Body))
	return nil
}

// mergeDocuments объединяет несколько PDF файлов
func mergeDocuments(client *gotenberg.Client) error {
	fmt.Println("\n=== Merging PDF Documents ===")

	// Сначала создаем два PDF документа для объединения
	pdf1Content := []byte(`<!DOCTYPE html>
		<html><body><h1>Document 1</h1><p>This is the first document.</p></body></html>`)

	pdf2Content := []byte(`<!DOCTYPE html>
		<html><body><h1>Document 2</h1><p>This is the second document.</p></body></html>`)

	// Создаем первый PDF
	resp1, err := client.ConvertHTMLToPDF(
		context.Background(),
		gotenberg.File{Name: "index.html", Content: pdf1Content},
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create PDF 1: %w", err)
	}

	// Создаем второй PDF
	resp2, err := client.ConvertHTMLToPDF(
		context.Background(),
		gotenberg.File{Name: "index.html", Content: pdf2Content},
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create PDF 2: %w", err)
	}

	// Объединяем PDF файлы
	mergedResp, err := client.MergePDFs(
		context.Background(),
		[]gotenberg.File{
			{Name: "doc1.pdf", Content: resp1.Body},
			{Name: "doc2.pdf", Content: resp2.Body},
		},
		&gotenberg.PDFEnginesRequest{
			OutputFilename: "merged",
		},
	)
	if err != nil {
		return fmt.Errorf("failed to merge PDFs: %w", err)
	}

	filename := "merged.pdf"
	if err := os.WriteFile(filename, mergedResp.Body, 0644); err != nil {
		return err
	}

	fmt.Printf("✓ Merged PDF created: %s (size: %d bytes)\n", filename, len(mergedResp.Body))
	return nil
}
