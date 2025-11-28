package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateTagHandler struct {
	documentTagService service.DocumentTagService
}

func NewUpdateTagHandler(documentTagService service.DocumentTagService) *UpdateTagHandler {
	return &UpdateTagHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Обновление тега
// @Description  Изменяет название тега
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID тега" format:"uuid"
// @Param        request body UpdateTagRequest true "Новое название"
// @Success      200 {object} TagResponse "Тег успешно обновлен"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Тег не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/tags/{id} [put]
func (h *UpdateTagHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	tagID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid tag id format",
		})
	}

	var req UpdateTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	tag, err := h.documentTagService.UpdateTag(c.Context(), tagID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(TagResponse{
		ID:        tag.ID,
		CompanyID: tag.CompanyID,
		Name:      tag.Name,
	})
}
