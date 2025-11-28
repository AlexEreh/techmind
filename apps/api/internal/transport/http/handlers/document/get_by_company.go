package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByCompanyHandler struct {
	documentService service.DocumentService
}

func NewGetByCompanyHandler(documentService service.DocumentService) *GetByCompanyHandler {
	return &GetByCompanyHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Получение всех документов компании
// @Description  Возвращает список всех документов компании
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company_id path string true "ID компании" format:"uuid"
// @Success      200 {object} DocumentsListResponse "Список документов"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/company/{company_id} [get]
func (h *GetByCompanyHandler) Handle(c fiber.Ctx) error {
	companyIDParam := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid company id format",
		})
	}

	docsWithTags, err := h.documentService.GetByCompany(c.Context(), companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := DocumentsListResponse{
		Documents: make([]DocumentResponse, 0, len(docsWithTags)),
		Total:     len(docsWithTags),
	}

	for _, docWithTags := range docsWithTags {
		tags := make([]TagData, 0, len(docWithTags.Tags))
		for _, tag := range docWithTags.Tags {
			tags = append(tags, TagData{
				ID:        tag.ID,
				CompanyID: tag.CompanyID,
				Name:      tag.Name,
			})
		}

		response.Documents = append(response.Documents, DocumentResponse{
			ID:              docWithTags.Document.ID,
			CompanyID:       docWithTags.Document.CompanyID,
			FolderID:        docWithTags.Document.FolderID,
			SenderID:        docWithTags.Document.SenderID,
			FilePath:        docWithTags.Document.FilePath,
			PreviewFilePath: docWithTags.Document.PreviewFilePath,
			FileSize:        docWithTags.Document.FileSize,
			MimeType:        docWithTags.Document.MimeType,
			Checksum:        docWithTags.Document.Checksum,
			CreatedAt:       docWithTags.Document.CreatedAt,
			Tags:            tags,
			PreviewURL:      docWithTags.PreviewURL,
		})
	}

	return c.JSON(response)
}
