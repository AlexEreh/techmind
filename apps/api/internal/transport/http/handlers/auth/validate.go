package auth

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type ValidateHandler struct {
	authService service.AuthService
}

func NewValidateHandler(authService service.AuthService) *ValidateHandler {
	return &ValidateHandler{
		authService: authService,
	}
}

// Handle godoc
// @Summary      Проверка токена
// @Description  Валидация JWT токена и получение ID пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body ValidateTokenRequest true "Токен для проверки"
// @Success      200 {object} ValidateTokenResponse "Токен валиден"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      401 {object} handlers.ErrorResponse "Токен невалиден"
// @Router       /public/auth/validate [post]
func (h *ValidateHandler) Handle(c fiber.Ctx) error {
	var req ValidateTokenRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	userID, err := h.authService.ValidateToken(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(ValidateTokenResponse{
			Valid: false,
		})
	}

	return c.JSON(ValidateTokenResponse{
		Valid:  true,
		UserID: userID,
	})
}
