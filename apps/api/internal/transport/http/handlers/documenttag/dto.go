package documenttag

import (
	"github.com/google/uuid"
)

// AddTagRequest представляет запрос на добавление тега к документу
type AddTagRequest struct {
	DocumentID uuid.UUID `json:"document_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	TagID      uuid.UUID `json:"tag_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// RemoveTagRequest представляет запрос на удаление тега у документа
type RemoveTagRequest struct {
	DocumentID uuid.UUID `json:"document_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	TagID      uuid.UUID `json:"tag_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// CreateTagRequest представляет запрос на создание тега
type CreateTagRequest struct {
	CompanyID uuid.UUID `json:"company_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" validate:"required,min=1" example:"Important"`
}

// UpdateTagRequest представляет запрос на обновление тега
type UpdateTagRequest struct {
	Name string `json:"name" validate:"required,min=1" example:"Very Important"`
}

// TagResponse представляет данные тега
type TagResponse struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID uuid.UUID `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name      string    `json:"name" example:"Important"`
}

// TagsListResponse представляет список тегов
type TagsListResponse struct {
	Tags  []TagResponse `json:"tags"`
	Total int           `json:"total" example:"5"`
}

// SuccessResponse представляет успешный ответ
type SuccessResponse struct {
	Message string `json:"message" example:"tag added successfully"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error" example:"tag not found"`
}
