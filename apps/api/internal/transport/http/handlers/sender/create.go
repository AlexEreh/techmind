package sender

import (
	"log"
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type CreateHandler struct {
	senderService service.SenderService
}

func NewCreateHandler(senderService service.SenderService) *CreateHandler {
	return &CreateHandler{
		senderService: senderService,
	}
}

// Handle godoc
// @Summary      Создание нового контрагента
// @Description  Создает нового контрагента в компании
// @Tags         senders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateSenderRequest true "Данные для создания контрагента"
// @Success      201 {object} SenderResponse "Контрагент успешно создан"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/senders [post]
func (h *CreateHandler) Handle(c fiber.Ctx) error {
	var req CreateSenderRequest
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

	sender, err := h.senderService.Create(c.Context(), req.CompanyID, req.Name, req.Email)
	if err != nil {
		log.Printf("Failed to create sender: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(SenderResponse{
		ID:        sender.ID,
		CompanyID: sender.CompanyID,
		Name:      sender.Name,
		Email:     sender.Email,
	})
}
