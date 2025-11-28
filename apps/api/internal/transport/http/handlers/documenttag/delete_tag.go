package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteTagHandler struct {
	documentTagService service.DocumentTagService
}

func NewDeleteTagHandler(documentTagService service.DocumentTagService) *DeleteTagHandler {
	return &DeleteTagHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Удаление тега
// @Description  Удаляет тег из системы и все его связи с документами
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID тега" format:"uuid"
// @Success      204 "Тег успешно удален"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Тег не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/tags/{id} [delete]
func (h *DeleteTagHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	tagID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid tag id format",
		})
	}

	if err := h.documentTagService.DeleteTag(c.Context(), tagID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
