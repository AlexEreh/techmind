package tag

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/tag"

	"github.com/google/uuid"
)

type tagRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.TagRepository {
	return &tagRepo{client: client}
}

func (r *tagRepo) Create(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error) {
	return r.client.Tag.
		Create().
		SetCompanyID(companyID).
		SetName(name).
		Save(ctx)
}

func (r *tagRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Tag, error) {
	return r.client.Tag.
		Query().
		Where(tag.ID(id)).
		Only(ctx)
}

func (r *tagRepo) GetByName(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error) {
	return r.client.Tag.
		Query().
		Where(
			tag.CompanyID(companyID),
			tag.Name(name),
		).
		Only(ctx)
}

func (r *tagRepo) Update(ctx context.Context, id uuid.UUID, name string) (*ent.Tag, error) {
	return r.client.Tag.
		UpdateOneID(id).
		SetName(name).
		Save(ctx)
}

func (r *tagRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Tag.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *tagRepo) List(ctx context.Context) ([]*ent.Tag, error) {
	return r.client.Tag.
		Query().
		All(ctx)
}

func (r *tagRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Tag, error) {
	return r.client.Tag.
		Query().
		Where(tag.CompanyID(companyID)).
		All(ctx)
}
