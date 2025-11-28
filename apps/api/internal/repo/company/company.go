package company

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/company"

	"github.com/google/uuid"
)

type companyRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.CompanyRepository {
	return &companyRepo{client: client}
}

func (r *companyRepo) Create(ctx context.Context, name string) (*ent.Company, error) {
	return r.client.Company.
		Create().
		SetName(name).
		Save(ctx)
}

func (r *companyRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Company, error) {
	return r.client.Company.
		Query().
		Where(company.ID(id)).
		Only(ctx)
}

func (r *companyRepo) Update(ctx context.Context, id uuid.UUID, name string) (*ent.Company, error) {
	return r.client.Company.
		UpdateOneID(id).
		SetName(name).
		Save(ctx)
}

func (r *companyRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Company.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *companyRepo) List(ctx context.Context) ([]*ent.Company, error) {
	return r.client.Company.
		Query().
		All(ctx)
}
