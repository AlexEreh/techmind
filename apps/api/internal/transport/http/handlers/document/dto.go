package document

import (
	"time"

	"github.com/google/uuid"
)

// UploadRequest представляет запрос на загрузку документа
type UploadRequest struct {
	CompanyID uuid.UUID  `form:"company_id" validate:"required"`
	FolderID  *uuid.UUID `form:"folder_id,omitempty"`
	Name      string     `form:"name" validate:"required"`
	SenderID  *uuid.UUID `form:"sender_id,omitempty"`
}

// DocumentResponse представляет данные документа
type DocumentResponse struct {
	ID              uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID       uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	FolderID        *uuid.UUID `json:"folder_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440002"`
	SenderID        *uuid.UUID `json:"sender_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440003"`
	Name            string     `json:"name" example:"document.pdf"`
	FilePath        string     `json:"file_path" example:"documents/550e8400-e29b-41d4-a716-446655440000.pdf"`
	PreviewFilePath *string    `json:"preview_file_path,omitempty" example:"previews/550e8400-e29b-41d4-a716-446655440000.jpg"`
	FileSize        int64      `json:"file_size" example:"1024000"`
	MimeType        string     `json:"mime_type" example:"application/pdf"`
	Checksum        string     `json:"checksum" example:"abc123def456"`
	CreatedAt       time.Time  `json:"created_at" example:"2024-11-28T15:04:05Z"`
	Tags            []TagData  `json:"tags,omitempty"`
	PreviewURL      string     `json:"preview_url,omitempty" example:"https://minio.example.com/bucket/preview.jpg?token=..."`
	DownloadURL     string     `json:"download_url,omitempty" example:"https://minio.example.com/bucket/document.pdf?token=..."`
}

// TagData представляет данные тега
type TagData struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	CompanyID uuid.UUID `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name      string    `json:"name" example:"Important"`
}

// UpdateRequest представляет запрос на обновление документа
type UpdateRequest struct {
	Name     string     `json:"name,omitempty" validate:"omitempty,min=1" example:"new_name.pdf"`
	FolderID *uuid.UUID `json:"folder_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	SenderID *uuid.UUID `json:"sender_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
}

// SearchRequest представляет запрос на поиск документов
type SearchRequest struct {
	CompanyID uuid.UUID   `json:"company_id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Query     string      `json:"query,omitempty" example:"invoice"`
	FolderID  *uuid.UUID  `json:"folder_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	TagIDs    []uuid.UUID `json:"tag_ids,omitempty"`
	Page      int         `json:"page,omitempty" example:"1"`
	PageSize  int         `json:"page_size,omitempty" example:"20"`
}

// DocumentsListResponse представляет список документов
type DocumentsListResponse struct {
	Documents []DocumentResponse `json:"documents"`
	Total     int                `json:"total" example:"10"`
}

// URLResponse представляет ответ с URL
type URLResponse struct {
	URL       string    `json:"url" example:"https://minio.example.com/bucket/document.pdf?token=..."`
	ExpiresAt time.Time `json:"expires_at" example:"2024-11-28T16:04:05Z"`
}
