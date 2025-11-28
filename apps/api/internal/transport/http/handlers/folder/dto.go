package folder

import "github.com/google/uuid"

// CreateRequest представляет запрос на создание папки
type CreateRequest struct {
	CompanyID uuid.UUID  `json:"company_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string     `json:"name" validate:"required,min=1" example:"Documents"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// FolderResponse представляет данные папки
type FolderResponse struct {
	ID             uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID      uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	ParentFolderID *uuid.UUID `json:"parent_folder_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440002"`
	Name           string     `json:"name" example:"Documents"`
	Size           int64      `json:"size" example:"0"`
	Count          int        `json:"count" example:"0"`
}

// RenameRequest представляет запрос на переименование папки
type RenameRequest struct {
	Name string `json:"name" validate:"required,min=1" example:"New Folder Name"`
}

// GetByParentRequest представляет запрос на получение вложенных папок
type GetByParentRequest struct {
	CompanyID uuid.UUID  `json:"company_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// FoldersListResponse представляет список папок
type FoldersListResponse struct {
	Folders []FolderResponse `json:"folders"`
	Total   int              `json:"total" example:"10"`
}
