package http

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const UserIDContextKey = "user_id"

// jwtMiddleware проверяет JWT токен в заголовке Authorization
func (s *Server) jwtMiddleware(c fiber.Ctx) error {
	// Получаем заголовок Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	// Извлекаем токен из заголовка "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization header format",
		})
	}

	tokenString := parts[1]

	// Валидируем токен
	userID, err := s.deps.AuthService.ValidateToken(context.Background(), tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or expired token",
		})
	}

	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid user id in token",
		})
	}

	// Сохраняем user_id в контексте для использования в handlers
	c.Locals(UserIDContextKey, userID)

	return c.Next()
}

// GetUserIDFromContext извлекает user_id из контекста
func GetUserIDFromContext(c fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals(UserIDContextKey)
	if userID == nil {
		return uuid.Nil, errors.New("user_id not found in context")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user_id type in context")
	}

	return id, nil
}
