package repo

import (
	"context"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

// UserRepository defines user-related database operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, name, email, password string) (*ent.User, error)
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error)
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*ent.User, error)
	// Update updates an existing user
	Update(ctx context.Context, id uuid.UUID, name, email, password string) (*ent.User, error)
	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all users
	List(ctx context.Context) ([]*ent.User, error)
}

// CompanyRepository defines company-related database operations
type CompanyRepository interface {
	// Create creates a new company
	Create(ctx context.Context, name string) (*ent.Company, error)
	// GetByID retrieves a company by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Company, error)
	// Update updates an existing company
	Update(ctx context.Context, id uuid.UUID, name string) (*ent.Company, error)
	// Delete deletes a company by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all companies
	List(ctx context.Context) ([]*ent.Company, error)
}

// CompanyUserRepository defines company user relationship operations
type CompanyUserRepository interface {
	// Create creates a new company user relationship
	Create(ctx context.Context, userID, companyID uuid.UUID, role int) (*ent.CompanyUser, error)
	// GetByID retrieves a company user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.CompanyUser, error)
	// GetByUserAndCompany retrieves a company user by user and company IDs
	GetByUserAndCompany(ctx context.Context, userID, companyID uuid.UUID) (*ent.CompanyUser, error)
	// GetUserRole retrieves the role of a user in a company
	GetUserRole(ctx context.Context, userID, companyID uuid.UUID) (int, error)
	// Update updates an existing company user relationship
	Update(ctx context.Context, id uuid.UUID, role int) (*ent.CompanyUser, error)
	// UpdateRole updates the role of a user in a company
	UpdateRole(ctx context.Context, userID, companyID uuid.UUID, newRole int) error
	// Delete deletes a company user relationship by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all company user relationships
	List(ctx context.Context) ([]*ent.CompanyUser, error)
	// ListByCompany retrieves all company users by company ID
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.CompanyUser, error)
	// ListByUser retrieves all company users by user ID
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*ent.CompanyUser, error)
}

// FolderRepository defines folder-related database operations
type FolderRepository interface {
	// Create creates a new folder
	Create(ctx context.Context, companyID uuid.UUID, parentFolderID *uuid.UUID, name string) (*ent.Folder, error)
	// GetByID retrieves a folder by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Folder, error)
	// Update updates an existing folder
	Update(ctx context.Context, id uuid.UUID, name string, size int64, count int) (*ent.Folder, error)
	// Delete deletes a folder by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all folders
	List(ctx context.Context) ([]*ent.Folder, error)
	// ListByCompany retrieves all folders for a company
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Folder, error)
	// ListByParent retrieves all child folders of a parent folder
	ListByParent(ctx context.Context, parentFolderID uuid.UUID) ([]*ent.Folder, error)
}

// SenderRepository defines sender-related database operations
type SenderRepository interface {
	// Create creates a new sender
	Create(ctx context.Context, name string, email *string) (*ent.Sender, error)
	// GetByID retrieves a sender by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Sender, error)
	// Update updates an existing sender
	Update(ctx context.Context, id uuid.UUID, name string, email *string) (*ent.Sender, error)
	// Delete deletes a sender by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all senders
	List(ctx context.Context) ([]*ent.Sender, error)
}

// DocumentRepository defines document-related database operations
type DocumentRepository interface {
	// Create creates a new document
	Create(ctx context.Context, companyID uuid.UUID, folderID *uuid.UUID, name string, filePath string, fileSize int64, mimeType string, checksum string) (*ent.Document, error)
	// GetByID retrieves a document by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Document, error)
	// Update updates an existing document
	Update(ctx context.Context, id uuid.UUID, folderID *uuid.UUID, senderID *uuid.UUID, name string) (*ent.Document, error)
	// UpdatePreviewPath updates the preview file path of a document
	UpdatePreviewPath(ctx context.Context, id uuid.UUID, previewFilePath string) error
	// Delete deletes a document by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all documents
	List(ctx context.Context) ([]*ent.Document, error)
	// ListByCompany retrieves all documents for a company
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Document, error)
	// ListByFolder retrieves all documents in a folder
	ListByFolder(ctx context.Context, folderID uuid.UUID) ([]*ent.Document, error)
}

// TagRepository defines tag-related database operations
type TagRepository interface {
	// Create creates a new tag for a company
	Create(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error)
	// GetByID retrieves a tag by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Tag, error)
	// GetByName retrieves a tag by name (within a company)
	GetByName(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error)
	// Update updates an existing tag
	Update(ctx context.Context, id uuid.UUID, name string) (*ent.Tag, error)
	// Delete deletes a tag by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all tags
	List(ctx context.Context) ([]*ent.Tag, error)
	// ListByCompany retrieves all tags for a company
	ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Tag, error)
}

// DocumentTagRepository defines document tag relationship operations
type DocumentTagRepository interface {
	// Create creates a new document tag relationship
	Create(ctx context.Context, documentID, tagID uuid.UUID) (*ent.DocumentTag, error)
	// GetByID retrieves a document tag by ID
	GetByID(ctx context.Context, id uuid.UUID) (*ent.DocumentTag, error)
	// Delete deletes a document tag relationship by ID
	Delete(ctx context.Context, id uuid.UUID) error
	// List retrieves all document tag relationships
	List(ctx context.Context) ([]*ent.DocumentTag, error)
	// ListByDocument retrieves all tags for a document
	ListByDocument(ctx context.Context, documentID uuid.UUID) ([]*ent.DocumentTag, error)
	// ListByTag retrieves all documents with a tag
	ListByTag(ctx context.Context, tagID uuid.UUID) ([]*ent.DocumentTag, error)
	// DeleteByDocumentAndTag deletes a document tag relationship
	DeleteByDocumentAndTag(ctx context.Context, documentID, tagID uuid.UUID) error
}
