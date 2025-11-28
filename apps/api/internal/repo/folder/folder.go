package folder

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/folder"

	"github.com/google/uuid"
)

type folderRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.FolderRepository {
	return &folderRepo{client: client}
}

func (r *folderRepo) Create(ctx context.Context, companyID uuid.UUID, parentFolderID *uuid.UUID, name string) (*ent.Folder, error) {
	create := r.client.Folder.
		Create().
		SetCompanyID(companyID).
		SetName(name)

	if parentFolderID != nil {
		create = create.SetParentFolderID(*parentFolderID)
	}

	return create.Save(ctx)
}

func (r *folderRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Folder, error) {
	return r.client.Folder.
		Query().
		Where(folder.ID(id)).
		Only(ctx)
}

func (r *folderRepo) Update(ctx context.Context, id uuid.UUID, name string, size int64, count int) (*ent.Folder, error) {
	return r.client.Folder.
		UpdateOneID(id).
		SetName(name).
		SetSize(size).
		SetCount(count).
		Save(ctx)
}

func (r *folderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Folder.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *folderRepo) List(ctx context.Context) ([]*ent.Folder, error) {
	return r.client.Folder.
		Query().
		All(ctx)
}

func (r *folderRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Folder, error) {
	return r.client.Folder.
		Query().
		Where(folder.CompanyID(companyID)).
		All(ctx)
}

func (r *folderRepo) ListByParent(ctx context.Context, parentFolderID uuid.UUID) ([]*ent.Folder, error) {
	return r.client.Folder.
		Query().
		Where(folder.ParentFolderID(parentFolderID)).
		All(ctx)
}
