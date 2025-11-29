package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByIDHandler struct {
	documentService service.DocumentService
}

func NewGetByIDHandler(documentService service.DocumentService) *GetByIDHandler {
	return &GetByIDHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Получение документа по ID
// @Description  Возвращает информацию о документе с тегами и preview URL
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID документа" format:"uuid"
// @Success      200 {object} DocumentResponse "Данные документа"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Документ не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/{id} [get]
func (h *GetByIDHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	documentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	docWithTags, err := h.documentService.GetByID(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(handlers.ErrorResponse{
			Error: "document not found",
		})
	}

	tags := make([]TagData, 0, len(docWithTags.Tags))
	for _, tag := range docWithTags.Tags {
		tags = append(tags, TagData{
			ID:        tag.ID,
			CompanyID: tag.CompanyID,
			Name:      tag.Name,
		})
	}

	return c.JSON(DocumentResponse{
		ID:              docWithTags.Document.ID,
		CompanyID:       docWithTags.Document.CompanyID,
		FolderID:        docWithTags.Document.FolderID,
		SenderID:        docWithTags.Document.SenderID,
		Name:            docWithTags.Document.Name,
		FilePath:        docWithTags.Document.FilePath,
		PreviewFilePath: docWithTags.Document.PreviewFilePath,
		FileSize:        docWithTags.Document.FileSize,
		MimeType:        docWithTags.Document.MimeType,
		Checksum:        docWithTags.Document.Checksum,
		CreatedBy:       docWithTags.Document.CreatedBy,
		UpdatedBy:       docWithTags.Document.UpdatedBy,
		CreatedAt:       docWithTags.Document.CreatedAt,
		UpdatedAt:       docWithTags.Document.UpdatedAt,
		Tags:            tags,
		PreviewURL:      docWithTags.PreviewURL,
		DownloadURL:     docWithTags.DownloadURL,
	})
}
