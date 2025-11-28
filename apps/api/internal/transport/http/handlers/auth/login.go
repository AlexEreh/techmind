package auth

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type LoginHandler struct {
	authService service.AuthService
}

func NewLoginHandler(authService service.AuthService) *LoginHandler {
	return &LoginHandler{
		authService: authService,
	}
}

// Handle godoc
// @Summary      Вход в систему
// @Description  Аутентификация пользователя по email и паролю
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Данные для входа"
// @Success      200 {object} LoginResponse "Успешная авторизация"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      401 {object} handlers.ErrorResponse "Неверные учетные данные"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /public/auth/login [post]
func (h *LoginHandler) Handle(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	token, expiresAt, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(handlers.ErrorResponse{
			Error: "invalid credentials",
		})
	}

	return c.JSON(LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}
