package gotenberg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
)

// PDFEnginesRequest представляет базовый запрос для PDF Engines операций
type PDFEnginesRequest struct {
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

	// Output Filename
	OutputFilename string
}

// SplitRequest представляет запрос для разделения PDF
type SplitRequest struct {
	PDFEnginesRequest
	SplitMode  string
	SplitSpan  string
	SplitUnify bool
}

// ConvertToPDFA конвертирует PDF файлы в PDF/A и/или PDF/UA
// POST /forms/pdfengines/convert
func (c *Client) ConvertToPDFA(ctx context.Context, files []File, pdfa string, pdfua bool, outputFilename string) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	if pdfa == "" && !pdfua {
		return nil, fmt.Errorf("at least one of pdfa or pdfua must be provided")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	fields := make(map[string]string)
	if pdfa != "" {
		fields["pdfa"] = pdfa
	}
	if pdfua {
		fields["pdfua"] = "true"
	}

	if err := addFieldsToWriter(writer, fields); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if outputFilename != "" {
		headers["Gotenberg-Output-Filename"] = outputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/convert", body, writer.FormDataContentType(), headers)
}

// ReadPDFMetadata читает метаданные PDF файлов
// POST /forms/pdfengines/metadata/read
func (c *Client) ReadPDFMetadata(ctx context.Context, files []File) (map[string]map[string]interface{}, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/forms/pdfengines/metadata/read", body, writer.FormDataContentType(), nil)
	if err != nil {
		return nil, err
	}

	var metadata map[string]map[string]interface{}
	if err := json.Unmarshal(resp.Body, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return metadata, nil
}

// WritePDFMetadata записывает метаданные в PDF файлы
// POST /forms/pdfengines/metadata/write
func (c *Client) WritePDFMetadata(ctx context.Context, files []File, metadata map[string]interface{}, outputFilename string) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	if len(metadata) == 0 {
		return nil, fmt.Errorf("metadata is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	// Добавляем метаданные
	metadataJSON, err := marshalJSONField(metadata)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("metadata", metadataJSON); err != nil {
		return nil, fmt.Errorf("failed to write metadata field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if outputFilename != "" {
		headers["Gotenberg-Output-Filename"] = outputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/metadata/write", body, writer.FormDataContentType(), headers)
}

// MergePDFs объединяет несколько PDF файлов в один
// POST /forms/pdfengines/merge
func (c *Client) MergePDFs(ctx context.Context, files []File, req *PDFEnginesRequest) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	// Добавляем поля запроса
	if err := addPDFEnginesFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/merge", body, writer.FormDataContentType(), headers)
}

// SplitPDFs разделяет PDF файлы
// POST /forms/pdfengines/split
func (c *Client) SplitPDFs(ctx context.Context, files []File, req *SplitRequest) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	if req == nil || req.SplitMode == "" || req.SplitSpan == "" {
		return nil, fmt.Errorf("splitMode and splitSpan are required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	// Добавляем обязательные поля
	if err := writer.WriteField("splitMode", req.SplitMode); err != nil {
		return nil, fmt.Errorf("failed to write splitMode field: %w", err)
	}
	if err := writer.WriteField("splitSpan", req.SplitSpan); err != nil {
		return nil, fmt.Errorf("failed to write splitSpan field: %w", err)
	}
	if req.SplitUnify {
		if err := writer.WriteField("splitUnify", "true"); err != nil {
			return nil, fmt.Errorf("failed to write splitUnify field: %w", err)
		}
	}

	// Добавляем остальные поля
	if err := addPDFEnginesFields(writer, &req.PDFEnginesRequest); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/split", body, writer.FormDataContentType(), headers)
}

// FlattenPDFs сглаживает PDF файлы
// POST /forms/pdfengines/flatten
func (c *Client) FlattenPDFs(ctx context.Context, files []File, outputFilename string) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if outputFilename != "" {
		headers["Gotenberg-Output-Filename"] = outputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/flatten", body, writer.FormDataContentType(), headers)
}

// EncryptPDFs добавляет защиту паролем к PDF файлам
// POST /forms/pdfengines/encrypt
func (c *Client) EncryptPDFs(ctx context.Context, files []File, userPassword, ownerPassword, outputFilename string) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	if userPassword == "" {
		return nil, fmt.Errorf("userPassword is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	// Добавляем пароли
	if err := writer.WriteField("userPassword", userPassword); err != nil {
		return nil, fmt.Errorf("failed to write userPassword field: %w", err)
	}
	if ownerPassword != "" {
		if err := writer.WriteField("ownerPassword", ownerPassword); err != nil {
			return nil, fmt.Errorf("failed to write ownerPassword field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if outputFilename != "" {
		headers["Gotenberg-Output-Filename"] = outputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/encrypt", body, writer.FormDataContentType(), headers)
}

// EmbedFilesInPDF встраивает файлы в PDF
// POST /forms/pdfengines/embed
func (c *Client) EmbedFilesInPDF(ctx context.Context, pdfFiles []File, embedFiles []File, outputFilename string) (*Response, error) {
	if len(pdfFiles) == 0 {
		return nil, fmt.Errorf("at least one PDF file is required")
	}

	if len(embedFiles) == 0 {
		return nil, fmt.Errorf("at least one file to embed is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем PDF файлы
	if err := addFilesToWriter(writer, pdfFiles, "files"); err != nil {
		return nil, err
	}

	// Добавляем файлы для встраивания
	if err := addFilesToWriter(writer, embedFiles, "embeds"); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if outputFilename != "" {
		headers["Gotenberg-Output-Filename"] = outputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/pdfengines/embed", body, writer.FormDataContentType(), headers)
}

// addPDFEnginesFields добавляет общие поля PDF Engines в multipart writer
func addPDFEnginesFields(writer *multipart.Writer, req *PDFEnginesRequest) error {
	if req == nil {
		return nil
	}

	fields := make(map[string]string)

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

	// Добавляем все поля
	if err := addFieldsToWriter(writer, fields); err != nil {
		return err
	}

	// Embed Files
	if len(req.Embeds) > 0 {
		if err := addFilesToWriter(writer, req.Embeds, "embeds"); err != nil {
			return err
		}
	}

	return nil
}
