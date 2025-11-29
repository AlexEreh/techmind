package sender

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByCompanyHandler struct {
	senderService service.SenderService
}

func NewGetByCompanyHandler(senderService service.SenderService) *GetByCompanyHandler {
	return &GetByCompanyHandler{
		senderService: senderService,
	}
}

// Handle godoc
// @Summary      Получение всех контрагентов компании
// @Description  Возвращает список всех контрагентов компании
// @Tags         senders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company_id path string true "ID компании" format:"uuid"
// @Success      200 {object} SendersListResponse "Список контрагентов"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/senders/company/{company_id} [get]
func (h *GetByCompanyHandler) Handle(c fiber.Ctx) error {
	companyIDParam := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid company id format",
		})
	}

	senders, err := h.senderService.GetByCompany(c.Context(), companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := SendersListResponse{
		Senders: make([]SenderResponse, 0, len(senders)),
		Total:   len(senders),
	}

	for _, sender := range senders {
		response.Senders = append(response.Senders, SenderResponse{
			ID:        sender.ID,
			CompanyID: sender.CompanyID,
			Name:      sender.Name,
			Email:     sender.Email,
		})
	}

	return c.JSON(response)
}
