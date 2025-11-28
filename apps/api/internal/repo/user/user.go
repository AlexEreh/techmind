package user

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/user"

	"github.com/google/uuid"
)

type userRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.UserRepository {
	return &userRepo{client: client}
}

func (r *userRepo) Create(ctx context.Context, name, email, password string) (*ent.User, error) {
	return r.client.User.
		Create().
		SetName(name).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.ID(id)).
		Only(ctx)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.
		Query().
		Where(user.Email(email)).
		Only(ctx)
}

func (r *userRepo) Update(ctx context.Context, id uuid.UUID, name, email, password string) (*ent.User, error) {
	return r.client.User.
		UpdateOneID(id).
		SetName(name).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.User.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *userRepo) List(ctx context.Context) ([]*ent.User, error) {
	return r.client.User.
		Query().
		All(ctx)
}
