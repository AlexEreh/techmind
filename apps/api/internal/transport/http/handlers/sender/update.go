package sender

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateHandler struct {
	senderService service.SenderService
}

func NewUpdateHandler(senderService service.SenderService) *UpdateHandler {
	return &UpdateHandler{
		senderService: senderService,
	}
}

// Handle godoc
// @Summary      Обновление контрагента
// @Description  Обновляет данные контрагента
// @Tags         senders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID контрагента" format:"uuid"
// @Param        request body UpdateSenderRequest true "Данные для обновления"
// @Success      200 {object} SenderResponse "Контрагент успешно обновлен"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Контрагент не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/senders/{id} [put]
func (h *UpdateHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid sender id format",
		})
	}

	var req UpdateSenderRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "name is required",
		})
	}

	sender, err := h.senderService.Update(c.Context(), id, req.Name, req.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(SenderResponse{
		ID:        sender.ID,
		CompanyID: sender.CompanyID,
		Name:      sender.Name,
		Email:     sender.Email,
	})
}
