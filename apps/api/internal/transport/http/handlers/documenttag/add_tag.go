package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type AddTagHandler struct {
	documentTagService service.DocumentTagService
}

func NewAddTagHandler(documentTagService service.DocumentTagService) *AddTagHandler {
	return &AddTagHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Добавление тега к документу
// @Description  Связывает существующий тег с документом
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body AddTagRequest true "ID документа и тега"
// @Success      200 {object} SuccessResponse "Тег успешно добавлен"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Документ или тег не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/add [post]
func (h *AddTagHandler) Handle(c fiber.Ctx) error {
	var req AddTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	if err := h.documentTagService.AddTagToDocument(c.Context(), req.DocumentID, req.TagID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "tag added successfully",
	})
}
