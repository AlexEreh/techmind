package auth

import (
	"time"

	"github.com/google/uuid"
)

// LoginRequest представляет запрос на вход в систему
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse представляет ответ при успешной авторизации
type LoginResponse struct {
	Token     string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresAt time.Time `json:"expires_at" example:"2024-11-30T15:04:05Z"`
}

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2" example:"John Doe"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// RegisterResponse представляет ответ при успешной регистрации
type RegisterResponse struct {
	Token string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserData `json:"user"`
}

// UserData представляет данные пользователя
type UserData struct {
	ID    uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string    `json:"name" example:"John Doe"`
	Email string    `json:"email" example:"user@example.com"`
}

// ValidateTokenRequest представляет запрос на проверку токена
type ValidateTokenRequest struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ValidateTokenResponse представляет ответ при проверке токена
type ValidateTokenResponse struct {
	Valid  bool      `json:"valid" example:"true"`
	UserID uuid.UUID `json:"user_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
}
