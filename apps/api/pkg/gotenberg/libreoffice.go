package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"strconv"
)

// LibreOfficeRequest представляет запрос для конвертации LibreOffice документов
type LibreOfficeRequest struct {
	// Page Properties
	Password                        string
	Landscape                       bool
	NativePageRanges                string
	UpdateIndexes                   bool
	ExportFormFields                bool
	AllowDuplicateFieldNames        bool
	ExportBookmarks                 bool
	ExportBookmarksToPDFDestination bool
	ExportPlaceholders              bool
	ExportNotes                     bool
	ExportNotesPages                bool
	ExportOnlyNotesPages            bool
	ExportNotesInMargin             bool
	ConvertOOOTargetToPDFTarget     bool
	ExportLinksRelativeFsys         bool
	ExportHiddenSlides              bool
	SkipEmptyPages                  bool
	AddOriginalDocumentAsStream     bool
	SinglePageSheets                bool

	// Compress
	LosslessImageCompression bool
	Quality                  int
	ReduceImageResolution    bool
	MaxImageResolution       int

	// Merge
	Merge bool

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

// ConvertOfficeToPDF конвертирует Office документы в PDF
// POST /forms/libreoffice/convert
func (c *Client) ConvertOfficeToPDF(ctx context.Context, files []File, req *LibreOfficeRequest) (*Response, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файлы
	if err := addFilesToWriter(writer, files, "files"); err != nil {
		return nil, err
	}

	// Добавляем поля запроса
	if err := addLibreOfficeFields(writer, req); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	headers := make(map[string]string)
	if req != nil && req.OutputFilename != "" {
		headers["Gotenberg-Output-Filename"] = req.OutputFilename
	}

	return c.doRequest(ctx, "POST", "/forms/libreoffice/convert", body, writer.FormDataContentType(), headers)
}

// addLibreOfficeFields добавляет поля LibreOffice запроса в multipart writer
func addLibreOfficeFields(writer *multipart.Writer, req *LibreOfficeRequest) error {
	if req == nil {
		return nil
	}

	fields := make(map[string]string)

	// Page Properties
	if req.Password != "" {
		fields["password"] = req.Password
	}
	if req.Landscape {
		fields["landscape"] = "true"
	}
	if req.NativePageRanges != "" {
		fields["nativePageRanges"] = req.NativePageRanges
	}
	if !req.UpdateIndexes {
		fields["updateIndexes"] = "false"
	}
	if !req.ExportFormFields {
		fields["exportFormFields"] = "false"
	}
	if req.AllowDuplicateFieldNames {
		fields["allowDuplicateFieldNames"] = "true"
	}
	if !req.ExportBookmarks {
		fields["exportBookmarks"] = "false"
	}
	if req.ExportBookmarksToPDFDestination {
		fields["exportBookmarksToPdfDestination"] = "true"
	}
	if req.ExportPlaceholders {
		fields["exportPlaceholders"] = "true"
	}
	if req.ExportNotes {
		fields["exportNotes"] = "true"
	}
	if req.ExportNotesPages {
		fields["exportNotesPages"] = "true"
	}
	if req.ExportOnlyNotesPages {
		fields["exportOnlyNotesPages"] = "true"
	}
	if req.ExportNotesInMargin {
		fields["exportNotesInMargin"] = "true"
	}
	if req.ConvertOOOTargetToPDFTarget {
		fields["convertOooTargetToPdfTarget"] = "true"
	}
	if req.ExportLinksRelativeFsys {
		fields["exportLinksRelativeFsys"] = "true"
	}
	if req.ExportHiddenSlides {
		fields["exportHiddenSlides"] = "true"
	}
	if req.SkipEmptyPages {
		fields["skipEmptyPages"] = "true"
	}
	if req.AddOriginalDocumentAsStream {
		fields["addOriginalDocumentAsStream"] = "true"
	}
	if req.SinglePageSheets {
		fields["singlePageSheets"] = "true"
	}

	// Compress
	if req.LosslessImageCompression {
		fields["losslessImageCompression"] = "true"
	}
	if req.Quality > 0 {
		fields["quality"] = strconv.Itoa(req.Quality)
	}
	if req.ReduceImageResolution {
		fields["reduceImageResolution"] = "true"
	}
	if req.MaxImageResolution > 0 {
		fields["maxImageResolution"] = strconv.Itoa(req.MaxImageResolution)
	}

	// Merge
	if req.Merge {
		fields["merge"] = "true"
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

	// Embed Files
	if len(req.Embeds) > 0 {
		if err := addFilesToWriter(writer, req.Embeds, "embeds"); err != nil {
			return err
		}
	}

	return nil
}
