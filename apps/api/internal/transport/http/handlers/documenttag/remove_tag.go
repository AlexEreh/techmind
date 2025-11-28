package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type RemoveTagHandler struct {
	documentTagService service.DocumentTagService
}

func NewRemoveTagHandler(documentTagService service.DocumentTagService) *RemoveTagHandler {
	return &RemoveTagHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Удаление тега у документа
// @Description  Удаляет связь между тегом и документом
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body RemoveTagRequest true "ID документа и тега"
// @Success      200 {object} SuccessResponse "Тег успешно удален"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Связь не найдена"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/remove [post]
func (h *RemoveTagHandler) Handle(c fiber.Ctx) error {
	var req RemoveTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	if err := h.documentTagService.RemoveTagFromDocument(c.Context(), req.DocumentID, req.TagID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "tag removed successfully",
	})
}
