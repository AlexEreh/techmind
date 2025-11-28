package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

const UserIDContextKey = "user_id"
const CompanyIDContextKey = "company_id"

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

// GetCompanyIDFromContext извлекает company_id из контекста
func GetCompanyIDFromContext(c fiber.Ctx) (uuid.UUID, error) {
	companyID := c.Locals(CompanyIDContextKey)
	if companyID == nil {
		return uuid.Nil, errors.New("company_id not found in context")
	}

	id, ok := companyID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid company_id type in context")
	}

	return id, nil
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
