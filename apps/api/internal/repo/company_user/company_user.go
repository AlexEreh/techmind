package company_user

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/companyuser"

	"github.com/google/uuid"
)

type companyUserRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.CompanyUserRepository {
	return &companyUserRepo{client: client}
}

func (r *companyUserRepo) Create(ctx context.Context, userID, companyID uuid.UUID, role int) (*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Create().
		SetUserID(userID).
		SetCompanyID(companyID).
		SetRole(role).
		Save(ctx)
}

func (r *companyUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Query().
		Where(companyuser.ID(id)).
		Only(ctx)
}

func (r *companyUserRepo) GetByUserAndCompany(ctx context.Context, userID, companyID uuid.UUID) (*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Query().
		Where(
			companyuser.UserID(userID),
			companyuser.CompanyID(companyID),
		).
		Only(ctx)
}

func (r *companyUserRepo) GetUserRole(ctx context.Context, userID, companyID uuid.UUID) (int, error) {
	cu, err := r.GetByUserAndCompany(ctx, userID, companyID)
	if err != nil {
		return 0, err
	}
	return cu.Role, nil
}

func (r *companyUserRepo) Update(ctx context.Context, id uuid.UUID, role int) (*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		UpdateOneID(id).
		SetRole(role).
		Save(ctx)
}

func (r *companyUserRepo) UpdateRole(ctx context.Context, userID, companyID uuid.UUID, newRole int) error {
	return r.client.CompanyUser.
		Update().
		Where(
			companyuser.UserID(userID),
			companyuser.CompanyID(companyID),
		).
		SetRole(newRole).
		Exec(ctx)
}

func (r *companyUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.CompanyUser.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *companyUserRepo) List(ctx context.Context) ([]*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Query().
		All(ctx)
}

func (r *companyUserRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Query().
		Where(companyuser.CompanyID(companyID)).
		All(ctx)
}

func (r *companyUserRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]*ent.CompanyUser, error) {
	return r.client.CompanyUser.
		Query().
		Where(companyuser.UserID(userID)).
		All(ctx)
}
