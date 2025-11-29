package document

import (
	"context"

	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/document"

	"github.com/google/uuid"
)

type documentRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.DocumentRepository {
	return &documentRepo{client: client}
}

func (r *documentRepo) Create(ctx context.Context, companyID uuid.UUID, folderID *uuid.UUID, name string, filePath string, fileSize int64, mimeType string, checksum string, createdBy uuid.UUID) (*ent.Document, error) {
	create := r.client.Document.
		Create().
		SetCompanyID(companyID).
		SetName(name).
		SetFilePath(filePath).
		SetFileSize(fileSize).
		SetMimeType(mimeType).
		SetChecksum(checksum).
		SetCreatedBy(createdBy).
		SetUpdatedBy(createdBy)

	if folderID != nil {
		create = create.SetFolderID(*folderID)
	}

	return create.Save(ctx)
}

func (r *documentRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Document, error) {
	return r.client.Document.
		Query().
		Where(document.ID(id)).
		WithSender().
		Only(ctx)
}

// func (r *documentRepo) Update(ctx context.Context, id uuid.UUID, filePath string, fileSize int64, mimeType string, checksum string) (*ent.Document, error) {
func (r *documentRepo) Update(ctx context.Context, id uuid.UUID, folderID *uuid.UUID, senderID *uuid.UUID, name string, updatedBy uuid.UUID) (*ent.Document, error) {
	update := r.client.Document.
		UpdateOneID(id).
		SetName(name).
		SetUpdatedBy(updatedBy)

	if folderID != nil {
		update = update.SetFolderID(*folderID)
	}
	if senderID != nil {
		update = update.SetSenderID(*senderID)
	}

	return update.Save(ctx)
}

func (r *documentRepo) UpdatePreviewPath(ctx context.Context, id uuid.UUID, previewFilePath string) error {
	return r.client.Document.
		UpdateOneID(id).
		SetPreviewFilePath(previewFilePath).
		Exec(ctx)
}

func (r *documentRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Document.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *documentRepo) List(ctx context.Context) ([]*ent.Document, error) {
	return r.client.Document.
		Query().
		All(ctx)
}

func (r *documentRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Document, error) {
	return r.client.Document.
		Query().
		Where(document.CompanyID(companyID)).
		WithSender().
		All(ctx)
}

func (r *documentRepo) ListByFolder(ctx context.Context, folderID uuid.UUID) ([]*ent.Document, error) {
	return r.client.Document.
		Query().
		Where(document.FolderID(folderID)).
		WithSender().
		All(ctx)
}
