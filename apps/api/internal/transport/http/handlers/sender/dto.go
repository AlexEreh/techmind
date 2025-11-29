package sender

import (
	"github.com/google/uuid"
)

// CreateSenderRequest представляет запрос на создание контрагента
type CreateSenderRequest struct {
	CompanyID uuid.UUID `json:"company_id" validate:"required"`
	Name      string    `json:"name" validate:"required,min=1"`
	Email     *string   `json:"email,omitempty"`
}

// UpdateSenderRequest представляет запрос на обновление контрагента
type UpdateSenderRequest struct {
	Name  string  `json:"name" validate:"required,min=1"`
	Email *string `json:"email,omitempty"`
}

// SenderResponse представляет данные контрагента
type SenderResponse struct {
	ID        uuid.UUID `json:"id"`
	CompanyID uuid.UUID `json:"company_id"`
	Name      string    `json:"name"`
	Email     *string   `json:"email,omitempty"`
}

// SendersListResponse представляет список контрагентов
type SendersListResponse struct {
	Senders []SenderResponse `json:"senders"`
	Total   int              `json:"total"`
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Message string `json:"message"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}
