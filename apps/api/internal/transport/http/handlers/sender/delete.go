package sender

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteHandler struct {
	senderService service.SenderService
}

func NewDeleteHandler(senderService service.SenderService) *DeleteHandler {
	return &DeleteHandler{
		senderService: senderService,
	}
}

// Handle godoc
// @Summary      Удаление контрагента
// @Description  Удаляет контрагента из системы
// @Tags         senders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID контрагента" format:"uuid"
// @Success      200 {object} SuccessResponse "Контрагент успешно удален"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/senders/{id} [delete]
func (h *DeleteHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid sender id format",
		})
	}

	if err := h.senderService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(SuccessResponse{
		Message: "sender deleted successfully",
	})
}
