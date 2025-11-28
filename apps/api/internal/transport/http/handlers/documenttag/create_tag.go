package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type CreateTagHandler struct {
	documentTagService service.DocumentTagService
}

func NewCreateTagHandler(documentTagService service.DocumentTagService) *CreateTagHandler {
	return &CreateTagHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Создание нового тега
// @Description  Создает новый тег в компании
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTagRequest true "Данные для создания тега"
// @Success      201 {object} TagResponse "Тег успешно создан"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/tags [post]
func (h *CreateTagHandler) Handle(c fiber.Ctx) error {
	var req CreateTagRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	tag, err := h.documentTagService.CreateTag(c.Context(), req.CompanyID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(TagResponse{
		ID:        tag.ID,
		CompanyID: tag.CompanyID,
		Name:      tag.Name,
	})
}
