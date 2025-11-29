package service

import (
	"context"
	"io"
	"time"

	"techmind/schema/ent"

	"github.com/google/uuid"
)

// AuthService определяет интерфейс для работы с авторизацией и аутентификацией
type AuthService interface {
	// Login выполняет вход пользователя в систему
	// Принимает email и пароль, возвращает JWT токен со сроком действия 72 часа
	Login(ctx context.Context, email, password string) (token string, expiresAt time.Time, err error)

	// Register регистрирует нового пользователя в системе
	// Создает пользователя и возвращает JWT токен для автоматического входа
	Register(ctx context.Context, name, email, password string) (token string, user *ent.User, err error)

	// ValidateToken проверяет валидность JWT токена
	// Возвращает ID пользователя если токен валиден
	ValidateToken(ctx context.Context, token string) (userID uuid.UUID, err error)
}

// FolderService определяет интерфейс для работы с папками
type FolderService interface {
	// Create создает новую папку в компании
	// Если parentID указан, папка создается как подпапка
	Create(ctx context.Context, companyID uuid.UUID, name string, parentID *uuid.UUID) (*ent.Folder, error)

	// Delete удаляет папку по ID
	// Также удаляет все вложенные папки и документы
	Delete(ctx context.Context, folderID uuid.UUID) error

	// Rename переименовывает папку
	Rename(ctx context.Context, folderID uuid.UUID, newName string) (*ent.Folder, error)

	// GetByCompany получает список всех папок в компании
	// Возвращает все папки без учета иерархии
	GetByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Folder, error)

	// GetByParent получает список папок по родительской папке
	// Если parentID nil, возвращает корневые папки компании
	GetByParent(ctx context.Context, companyID uuid.UUID, parentID *uuid.UUID) ([]*ent.Folder, error)

	// GetByID получает папку по ID
	GetByID(ctx context.Context, folderID uuid.UUID) (*ent.Folder, error)
}

// DocumentUploadInput содержит данные для загрузки документа
type DocumentUploadInput struct {
	CompanyID uuid.UUID
	FolderID  *uuid.UUID
	Name      string
	File      io.Reader
	FileSize  int64
	MimeType  string
	SenderID  *uuid.UUID
}

// DocumentUpdateInput содержит данные для обновления метаданных документа
type DocumentUpdateInput struct {
	Name     string
	FolderID *uuid.UUID
	SenderID *uuid.UUID
}

// DocumentWithTags содержит документ вместе с его тегами
type DocumentWithTags struct {
	Document    *ent.Document
	Tags        []*ent.Tag
	PreviewURL  string
	DownloadURL string
}

// DocumentService определяет интерфейс для работы с документами
type DocumentService interface {
	// Upload загружает новый документ в систему
	// Принимает файл, сохраняет его в MinIO и создает запись в БД
	// Генерирует preview для поддерживаемых типов файлов
	Upload(ctx context.Context, input DocumentUploadInput) (*ent.Document, error)

	// GetByID получает документ по ID вместе с его тегами
	// Возвращает ссылку на preview документа
	GetByID(ctx context.Context, documentID uuid.UUID) (*DocumentWithTags, error)

	// GetByFolder получает список документов в папке
	// Включает теги для каждого документа
	GetByFolder(ctx context.Context, folderID uuid.UUID) ([]*DocumentWithTags, error)

	// GetByCompany получает список всех документов компании
	GetByCompany(ctx context.Context, companyID uuid.UUID) ([]*DocumentWithTags, error)

	// Update обновляет метаданные документа
	// Позволяет изменить имя, папку и отправителя
	Update(ctx context.Context, documentID uuid.UUID, input DocumentUpdateInput) (*ent.Document, error)

	// Delete удаляет документ из системы
	// Удаляет файлы из MinIO и запись из БД
	Delete(ctx context.Context, documentID uuid.UUID) error

	// GetDownloadURL получает временную ссылку на скачивание оригинала документа
	// Возвращает presigned URL для доступа к файлу в MinIO
	GetDownloadURL(ctx context.Context, documentID uuid.UUID) (url string, err error)

	// GetPreviewURL получает временную ссылку на preview документа
	// Возвращает presigned URL для доступа к preview файлу
	GetPreviewURL(ctx context.Context, documentID uuid.UUID) (url string, err error)

	// Search ищет документы по различным критериям
	Search(ctx context.Context, companyID uuid.UUID, query string, folderID *uuid.UUID, tagIDs []uuid.UUID) ([]*DocumentWithTags, error)

	// GeneratePDFPreview конвертирует файл документа в PDF превью и загружает его в MinIO
	// Поддерживает конвертацию Office документов (docx, xlsx, pptx и т.д.) через Gotenberg
	// После успешной конвертации обновляет ссылку на preview в базе данных
	GeneratePDFPreview(ctx context.Context, documentID uuid.UUID) error

	// ExtractAndIndexText извлекает текст из документа и индексирует его в Elasticsearch
	// Использует docconv для извлечения текста из различных форматов документов
	// Сохраняет извлеченный текст в индекс "documents" в Elasticsearch
	ExtractAndIndexText(ctx context.Context, documentID uuid.UUID) error
}

// DocumentTagService определяет интерфейс для работы с тегами документов
type DocumentTagService interface {
	// GetDocumentTags получает все теги конкретного документа
	GetDocumentTags(ctx context.Context, documentID uuid.UUID) ([]*ent.Tag, error)

	// AddTagToDocument добавляет существующий тег к документу
	AddTagToDocument(ctx context.Context, documentID, tagID uuid.UUID) error

	// RemoveTagFromDocument удаляет тег у документа
	RemoveTagFromDocument(ctx context.Context, documentID, tagID uuid.UUID) error

	// CreateTag создает новый тег в компании
	CreateTag(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error)

	// DeleteTag удаляет тег из системы
	// Также удаляет все связи этого тега с документами
	DeleteTag(ctx context.Context, tagID uuid.UUID) error

	// GetTagsByCompany получает все теги компании
	GetTagsByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Tag, error)

	// GetTagByID получает тег по ID
	GetTagByID(ctx context.Context, tagID uuid.UUID) (*ent.Tag, error)

	// UpdateTag обновляет название тега
	UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*ent.Tag, error)
}

// CompanyUserService определяет интерфейс для работы с пользователями компании
type CompanyUserService interface {
	// GetUserRole получает роль пользователя в конкретной компании
	GetUserRole(ctx context.Context, userID, companyID uuid.UUID) (int, error)

	// GetUserCompanies получает список всех компаний пользователя с информацией о ролях
	GetUserCompanies(ctx context.Context, userID uuid.UUID) ([]*ent.CompanyUser, error)
}

// CompanyService определяет интерфейс для работы с компаниями
type CompanyService interface {
	// Create создает новую компанию и добавляет создателя как администратора
	Create(ctx context.Context, name string, userID uuid.UUID) (*ent.Company, error)
}
