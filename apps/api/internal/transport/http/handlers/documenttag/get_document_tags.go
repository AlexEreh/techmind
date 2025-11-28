package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetDocumentTagsHandler struct {
	documentTagService service.DocumentTagService
}

func NewGetDocumentTagsHandler(documentTagService service.DocumentTagService) *GetDocumentTagsHandler {
	return &GetDocumentTagsHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Получение тегов документа
// @Description  Возвращает список всех тегов конкретного документа
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        document_id path string true "ID документа" format:"uuid"
// @Success      200 {object} TagsListResponse "Список тегов"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/document/{document_id} [get]
func (h *GetDocumentTagsHandler) Handle(c fiber.Ctx) error {
	documentIDParam := c.Params("document_id")
	documentID, err := uuid.Parse(documentIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	tags, err := h.documentTagService.GetDocumentTags(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := TagsListResponse{
		Tags:  make([]TagResponse, 0, len(tags)),
		Total: len(tags),
	}

	for _, tag := range tags {
		response.Tags = append(response.Tags, TagResponse{
			ID:        tag.ID,
			CompanyID: tag.CompanyID,
			Name:      tag.Name,
		})
	}

	return c.JSON(response)
}
