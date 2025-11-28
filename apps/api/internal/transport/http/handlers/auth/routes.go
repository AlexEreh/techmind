package auth

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с аутентификацией
func RegisterRoutes(router fiber.Router, authService service.AuthService) {
	loginHandler := NewLoginHandler(authService)
	registerHandler := NewRegisterHandler(authService)
	validateHandler := NewValidateHandler(authService)

	router.Post("/login", loginHandler.Handle)
	router.Post("/register", registerHandler.Handle)
	router.Post("/validate", validateHandler.Handle)
}
