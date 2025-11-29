package company_user

import (
	"time"

	"github.com/google/uuid"
)

// CompanyData содержит информацию о компании
type CompanyData struct {
	ID   uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Name string    `json:"name" example:"My Company"`
}

// CompanyUserData содержит информацию о связи пользователя и компании
type CompanyUserData struct {
	ID        uuid.UUID    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    uuid.UUID    `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	CompanyID uuid.UUID    `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Role      int          `json:"role" example:"1"`
	Company   *CompanyData `json:"company,omitempty"`
}

// MyCompaniesResponse содержит список компаний пользователя
type MyCompaniesResponse struct {
	Companies []CompanyUserData `json:"companies"`
}

// CompanyUserWithDetailsDTO содержит информацию о пользователе компании с деталями
type CompanyUserWithDetailsDTO struct {
	ID       uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Username string    `json:"username" example:"John Doe"`
	Email    string    `json:"email" example:"john@example.com"`
	Role     int       `json:"role" example:"1"`
	AddedAt  time.Time `json:"added_at" example:"2023-01-01T00:00:00Z"`
}
