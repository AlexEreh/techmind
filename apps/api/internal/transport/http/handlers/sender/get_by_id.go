package sender

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByIDHandler struct {
	senderService service.SenderService
}

func NewGetByIDHandler(senderService service.SenderService) *GetByIDHandler {
	return &GetByIDHandler{
		senderService: senderService,
	}
}

// Handle godoc
// @Summary      Получение контрагента по ID
// @Description  Возвращает данные контрагента по его ID
// @Tags         senders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID контрагента" format:"uuid"
// @Success      200 {object} SenderResponse "Данные контрагента"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Контрагент не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/senders/{id} [get]
func (h *GetByIDHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid sender id format",
		})
	}

	sender, err := h.senderService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(handlers.ErrorResponse{
			Error: "sender not found",
		})
	}

	return c.JSON(SenderResponse{
		ID:        sender.ID,
		CompanyID: sender.CompanyID,
		Name:      sender.Name,
		Email:     sender.Email,
	})
}
