package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetTagByIDHandler struct {
	documentTagService service.DocumentTagService
}

func NewGetTagByIDHandler(documentTagService service.DocumentTagService) *GetTagByIDHandler {
	return &GetTagByIDHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Получение тега по ID
// @Description  Возвращает информацию о теге по его идентификатору
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID тега" format:"uuid"
// @Success      200 {object} TagResponse "Данные тега"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Тег не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/tags/{id} [get]
func (h *GetTagByIDHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	tagID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid tag id format",
		})
	}

	tag, err := h.documentTagService.GetTagByID(c.Context(), tagID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(handlers.ErrorResponse{
			Error: "tag not found",
		})
	}

	return c.JSON(TagResponse{
		ID:        tag.ID,
		CompanyID: tag.CompanyID,
		Name:      tag.Name,
	})
}
