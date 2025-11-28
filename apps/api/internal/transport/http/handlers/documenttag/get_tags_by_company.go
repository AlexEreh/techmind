package documenttag

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetTagsByCompanyHandler struct {
	documentTagService service.DocumentTagService
}

func NewGetTagsByCompanyHandler(documentTagService service.DocumentTagService) *GetTagsByCompanyHandler {
	return &GetTagsByCompanyHandler{
		documentTagService: documentTagService,
	}
}

// Handle godoc
// @Summary      Получение всех тегов компании
// @Description  Возвращает список всех тегов компании
// @Tags         document-tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company_id path string true "ID компании" format:"uuid"
// @Success      200 {object} TagsListResponse "Список тегов"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/document-tags/company/{company_id} [get]
func (h *GetTagsByCompanyHandler) Handle(c fiber.Ctx) error {
	companyIDParam := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid company id format",
		})
	}

	tags, err := h.documentTagService.GetTagsByCompany(c.Context(), companyID)
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
