package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
)

// ChromiumRequest представляет базовый запрос для Chromium конвертации
type ChromiumRequest struct {
	// Page Properties
	SinglePage              bool
	PaperWidth              string
	PaperHeight             string
	MarginTop               string
	MarginBottom            string
	MarginLeft              string
	MarginRight             string
	PreferCSSPageSize       bool
	GenerateDocumentOutline bool
	GenerateTaggedPDF       bool
	PrintBackground         bool
	OmitBackground          bool
	Landscape               bool
	Scale                   float64
	NativePageRanges        string

	// Header & Footer
	Header *File
	Footer *File

	// Wait Before Rendering
	WaitDelay         string
	WaitForExpression string

	// Emulated Media Type
	EmulatedMediaType string

	// Cookies
	Cookies []Cookie

	// Custom HTTP Headers
	UserAgent        string
	ExtraHTTPHeaders map[string]string

	// Invalid HTTP Status Codes
	FailOnHTTPStatusCodes         []int
	FailOnResourceHTTPStatusCodes []int

	// Network Errors
	FailOnResourceLoadingFailed bool

	// Console Exceptions
	FailOnConsoleExceptions bool

	// Performance Mode
	SkipNetworkIdleEvent bool

	// Split
	SplitMode  string
	SplitSpan  string
	SplitUnify bool

	// PDF/A & PDF/UA
	PDFA  string
	PDFUA bool

	// Metadata
	Metadata map[string]interface{}

	// Flatten
	Flatten bool

	// Encrypt
	UserPassword  string
	OwnerPassword string

	// Embed Files
	Embeds []File

	// Download From
	DownloadFrom []DownloadFrom

	// Output Filename
	OutputFilename string
}

// Cookie представляет cookie для Chromium
type Cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path,omitempty"`
	Secure   bool   `json:"secure,omitempty"`
	HTTPOnly bool   `json:"httpOnly,omitempty"`
	SameSite string `json:"sameSite,omitempty"`
}

// ConvertURLToPDF конвертирует веб-страницу в PDF
// POST /forms/chromium/convert/url
func (c *Client) ConvertURLToPDF(ctx context.Context, url string, req *ChromiumRequest) (*Response, error) {
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем обязательное поле url
	if err := writer.WriteField("url", url); err != nil {
		return nil, fmt.Errorf("failed to write url field: %w", err)
	}

	// Добавляем остальные поля запроса
	if err := addChromiumFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/convert/url", body, writer.FormDataContentType(), headers)
}

// ConvertHTMLToPDF конвертирует HTML файл в PDF
// POST /forms/chromium/convert/html
func (c *Client) ConvertHTMLToPDF(ctx context.Context, indexHTML File, additionalFiles []File, req *ChromiumRequest) (*Response, error) {
	if indexHTML.Name == "" {
		indexHTML.Name = "index.html"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем index.html
	if err := addFilesToWriter(writer, []File{indexHTML}, "files"); err != nil {
		return nil, err
	}

	// Добавляем дополнительные файлы (изображения, стили и т.д.)
	if len(additionalFiles) > 0 {
		if err := addFilesToWriter(writer, additionalFiles, "files"); err != nil {
			return nil, err
		}
	}

	// Добавляем остальные поля запроса
	if err := addChromiumFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/convert/html", body, writer.FormDataContentType(), headers)
}

// ConvertMarkdownToPDF конвертирует Markdown файлы в PDF
// POST /forms/chromium/convert/markdown
func (c *Client) ConvertMarkdownToPDF(ctx context.Context, indexHTML File, markdownFiles []File, additionalFiles []File, req *ChromiumRequest) (*Response, error) {
	if indexHTML.Name == "" {
		indexHTML.Name = "index.html"
	}

	if len(markdownFiles) == 0 {
		return nil, fmt.Errorf("at least one markdown file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем index.html
	if err := addFilesToWriter(writer, []File{indexHTML}, "files"); err != nil {
		return nil, err
	}

	// Добавляем markdown файлы
	if err := addFilesToWriter(writer, markdownFiles, "files"); err != nil {
		return nil, err
	}

	// Добавляем дополнительные файлы
	if len(additionalFiles) > 0 {
		if err := addFilesToWriter(writer, additionalFiles, "files"); err != nil {
			return nil, err
		}
	}

	// Добавляем остальные поля запроса
	if err := addChromiumFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/convert/markdown", body, writer.FormDataContentType(), headers)
}

// ScreenshotRequest представляет запрос для создания скриншота
type ScreenshotRequest struct {
	// Screenshot Properties
	Width            int
	Height           int
	Clip             bool
	Format           string // png, jpeg, webp
	Quality          int    // 0-100, только для jpeg
	OmitBackground   bool
	OptimizeForSpeed bool

	// Wait Before Rendering
	WaitDelay         string
	WaitForExpression string

	// Emulated Media Type
	EmulatedMediaType string

	// Cookies
	Cookies []Cookie

	// Custom HTTP Headers
	UserAgent        string
	ExtraHTTPHeaders map[string]string

	// Invalid HTTP Status Codes
	FailOnHTTPStatusCodes         []int
	FailOnResourceHTTPStatusCodes []int

	// Console Exceptions
	FailOnConsoleExceptions bool

	// Performance Mode
	SkipNetworkIdleEvent bool

	// Output Filename
	OutputFilename string
}

// ScreenshotURL создает скриншот веб-страницы
// POST /forms/chromium/screenshot/url
func (c *Client) ScreenshotURL(ctx context.Context, url string, req *ScreenshotRequest) (*Response, error) {
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("url", url); err != nil {
		return nil, fmt.Errorf("failed to write url field: %w", err)
	}

	if err := addScreenshotFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/screenshot/url", body, writer.FormDataContentType(), headers)
}

// ScreenshotHTML создает скриншот HTML файла
// POST /forms/chromium/screenshot/html
func (c *Client) ScreenshotHTML(ctx context.Context, indexHTML File, additionalFiles []File, req *ScreenshotRequest) (*Response, error) {
	if indexHTML.Name == "" {
		indexHTML.Name = "index.html"
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := addFilesToWriter(writer, []File{indexHTML}, "files"); err != nil {
		return nil, err
	}

	if len(additionalFiles) > 0 {
		if err := addFilesToWriter(writer, additionalFiles, "files"); err != nil {
			return nil, err
		}
	}

	if err := addScreenshotFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/screenshot/html", body, writer.FormDataContentType(), headers)
}

// ScreenshotMarkdown создает скриншот Markdown файлов
// POST /forms/chromium/screenshot/markdown
func (c *Client) ScreenshotMarkdown(ctx context.Context, indexHTML File, markdownFiles []File, additionalFiles []File, req *ScreenshotRequest) (*Response, error) {
	if indexHTML.Name == "" {
		indexHTML.Name = "index.html"
	}

	if len(markdownFiles) == 0 {
		return nil, fmt.Errorf("at least one markdown file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := addFilesToWriter(writer, []File{indexHTML}, "files"); err != nil {
		return nil, err
	}

	if err := addFilesToWriter(writer, markdownFiles, "files"); err != nil {
		return nil, err
	}

	if len(additionalFiles) > 0 {
		if err := addFilesToWriter(writer, additionalFiles, "files"); err != nil {
			return nil, err
		}
	}

	if err := addScreenshotFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/chromium/screenshot/markdown", body, writer.FormDataContentType(), headers)
}

// addChromiumFields добавляет поля Chromium запроса в multipart writer
func addChromiumFields(writer *multipart.Writer, req *ChromiumRequest) error {
	if req == nil {
		return nil
	}

	fields := make(map[string]string)

	// Page Properties
	if req.SinglePage {
		fields["singlePage"] = "true"
	}
	if req.PaperWidth != "" {
		fields["paperWidth"] = req.PaperWidth
	}
	if req.PaperHeight != "" {
		fields["paperHeight"] = req.PaperHeight
	}
	if req.MarginTop != "" {
		fields["marginTop"] = req.MarginTop
	}
	if req.MarginBottom != "" {
		fields["marginBottom"] = req.MarginBottom
	}
	if req.MarginLeft != "" {
		fields["marginLeft"] = req.MarginLeft
	}
	if req.MarginRight != "" {
		fields["marginRight"] = req.MarginRight
	}
	if req.PreferCSSPageSize {
		fields["preferCssPageSize"] = "true"
	}
	if req.GenerateDocumentOutline {
		fields["generateDocumentOutline"] = "true"
	}
	if req.GenerateTaggedPDF {
		fields["generateTaggedPdf"] = "true"
	}
	if req.PrintBackground {
		fields["printBackground"] = "true"
	}
	if req.OmitBackground {
		fields["omitBackground"] = "true"
	}
	if req.Landscape {
		fields["landscape"] = "true"
	}
	if req.Scale > 0 {
		fields["scale"] = strconv.FormatFloat(req.Scale, 'f', -1, 64)
	}
	if req.NativePageRanges != "" {
		fields["nativePageRanges"] = req.NativePageRanges
	}

	// Wait Before Rendering
	if req.WaitDelay != "" {
		fields["waitDelay"] = req.WaitDelay
	}
	if req.WaitForExpression != "" {
		fields["waitForExpression"] = req.WaitForExpression
	}

	// Emulated Media Type
	if req.EmulatedMediaType != "" {
		fields["emulatedMediaType"] = req.EmulatedMediaType
	}

	// Cookies
	if len(req.Cookies) > 0 {
		cookiesJSON, err := marshalJSONField(req.Cookies)
		if err != nil {
			return err
		}
		fields["cookies"] = cookiesJSON
	}

	// Custom HTTP Headers
	if req.UserAgent != "" {
		fields["userAgent"] = req.UserAgent
	}
	if len(req.ExtraHTTPHeaders) > 0 {
		headersJSON, err := marshalJSONField(req.ExtraHTTPHeaders)
		if err != nil {
			return err
		}
		fields["extraHttpHeaders"] = headersJSON
	}

	// Invalid HTTP Status Codes
	if len(req.FailOnHTTPStatusCodes) > 0 {
		codesJSON, err := marshalJSONField(req.FailOnHTTPStatusCodes)
		if err != nil {
			return err
		}
		fields["failOnHttpStatusCodes"] = codesJSON
	}
	if len(req.FailOnResourceHTTPStatusCodes) > 0 {
		codesJSON, err := marshalJSONField(req.FailOnResourceHTTPStatusCodes)
		if err != nil {
			return err
		}
		fields["failOnResourceHttpStatusCodes"] = codesJSON
	}

	// Network Errors
	if req.FailOnResourceLoadingFailed {
		fields["failOnResourceLoadingFailed"] = "true"
	}

	// Console Exceptions
	if req.FailOnConsoleExceptions {
		fields["failOnConsoleExceptions"] = "true"
	}

	// Performance Mode
	if req.SkipNetworkIdleEvent {
		fields["skipNetworkIdleEvent"] = "true"
	}

	// Split
	if req.SplitMode != "" {
		fields["splitMode"] = req.SplitMode
	}
	if req.SplitSpan != "" {
		fields["splitSpan"] = req.SplitSpan
	}
	if req.SplitUnify {
		fields["splitUnify"] = "true"
	}

	// PDF/A & PDF/UA
	if req.PDFA != "" {
		fields["pdfa"] = req.PDFA
	}
	if req.PDFUA {
		fields["pdfua"] = "true"
	}

	// Metadata
	if len(req.Metadata) > 0 {
		metadataJSON, err := marshalJSONField(req.Metadata)
		if err != nil {
			return err
		}
		fields["metadata"] = metadataJSON
	}

	// Flatten
	if req.Flatten {
		fields["flatten"] = "true"
	}

	// Encrypt
	if req.UserPassword != "" {
		fields["userPassword"] = req.UserPassword
	}
	if req.OwnerPassword != "" {
		fields["ownerPassword"] = req.OwnerPassword
	}

	// Download From
	if len(req.DownloadFrom) > 0 {
		downloadJSON, err := marshalJSONField(req.DownloadFrom)
		if err != nil {
			return err
		}
		fields["downloadFrom"] = downloadJSON
	}

	// Добавляем все поля
	if err := addFieldsToWriter(writer, fields); err != nil {
		return err
	}

	// Header & Footer
	if req.Header != nil {
		if req.Header.Name == "" {
			req.Header.Name = "header.html"
		}
		if err := addFilesToWriter(writer, []File{*req.Header}, "files"); err != nil {
			return err
		}
	}
	if req.Footer != nil {
		if req.Footer.Name == "" {
			req.Footer.Name = "footer.html"
		}
		if err := addFilesToWriter(writer, []File{*req.Footer}, "files"); err != nil {
			return err
		}
	}

	// Embed Files
	if len(req.Embeds) > 0 {
		if err := addFilesToWriter(writer, req.Embeds, "embeds"); err != nil {
			return err
		}
	}

	return nil
}

// addScreenshotFields добавляет поля Screenshot запроса в multipart writer
func addScreenshotFields(writer *multipart.Writer, req *ScreenshotRequest) error {
	if req == nil {
		return nil
	}

	fields := make(map[string]string)

	// Screenshot Properties
	if req.Width > 0 {
		fields["width"] = strconv.Itoa(req.Width)
	}
	if req.Height > 0 {
		fields["height"] = strconv.Itoa(req.Height)
	}
	if req.Clip {
		fields["clip"] = "true"
	}
	if req.Format != "" {
		fields["format"] = req.Format
	}
	if req.Quality > 0 {
		fields["quality"] = strconv.Itoa(req.Quality)
	}
	if req.OmitBackground {
		fields["omitBackground"] = "true"
	}
	if req.OptimizeForSpeed {
		fields["optimizeForSpeed"] = "true"
	}

	// Wait Before Rendering
	if req.WaitDelay != "" {
		fields["waitDelay"] = req.WaitDelay
	}
	if req.WaitForExpression != "" {
		fields["waitForExpression"] = req.WaitForExpression
	}

	// Emulated Media Type
	if req.EmulatedMediaType != "" {
		fields["emulatedMediaType"] = req.EmulatedMediaType
	}

	// Cookies
	if len(req.Cookies) > 0 {
		cookiesJSON, err := marshalJSONField(req.Cookies)
		if err != nil {
			return err
		}
		fields["cookies"] = cookiesJSON
	}

	// Custom HTTP Headers
	if req.UserAgent != "" {
		fields["userAgent"] = req.UserAgent
	}
	if len(req.ExtraHTTPHeaders) > 0 {
		headersJSON, err := marshalJSONField(req.ExtraHTTPHeaders)
		if err != nil {
			return err
		}
		fields["extraHttpHeaders"] = headersJSON
	}

	// Invalid HTTP Status Codes
	if len(req.FailOnHTTPStatusCodes) > 0 {
		codesJSON, err := marshalJSONField(req.FailOnHTTPStatusCodes)
		if err != nil {
			return err
		}
		fields["failOnHttpStatusCodes"] = codesJSON
	}
	if len(req.FailOnResourceHTTPStatusCodes) > 0 {
		codesJSON, err := marshalJSONField(req.FailOnResourceHTTPStatusCodes)
		if err != nil {
			return err
		}
		fields["failOnResourceHttpStatusCodes"] = codesJSON
	}

	// Console Exceptions
	if req.FailOnConsoleExceptions {
		fields["failOnConsoleExceptions"] = "true"
	}

	// Performance Mode
	if req.SkipNetworkIdleEvent {
		fields["skipNetworkIdleEvent"] = "true"
	}

	return addFieldsToWriter(writer, fields)
}
