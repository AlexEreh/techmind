package auth

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type RegisterHandler struct {
	authService service.AuthService
}

func NewRegisterHandler(authService service.AuthService) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
	}
}

// Handle godoc
// @Summary      Регистрация нового пользователя
// @Description  Создание нового аккаунта пользователя в системе
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Данные для регистрации"
// @Success      201 {object} RegisterResponse "Успешная регистрация"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      409 {object} handlers.ErrorResponse "Пользователь уже существует"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /public/auth/register [post]
func (h *RegisterHandler) Handle(c fiber.Ctx) error {
	var req RegisterRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	token, user, err := h.authService.Register(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(handlers.ErrorResponse{
			Error: "user already exists or registration failed",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Token: token,
		User: UserData{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}
