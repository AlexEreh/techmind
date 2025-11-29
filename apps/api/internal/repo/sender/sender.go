package sender

import (
	"context"

	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/sender"

	"github.com/google/uuid"
)

type senderRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.SenderRepository {
	return &senderRepo{client: client}
}

func (r *senderRepo) Create(ctx context.Context, companyID uuid.UUID, name string, email *string) (*ent.Sender, error) {
	create := r.client.Sender.
		Create().
		SetCompanyID(companyID).
		SetName(name)

	if email != nil {
		create = create.SetEmail(*email)
	}

	return create.Save(ctx)
}

func (r *senderRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Sender, error) {
	return r.client.Sender.
		Query().
		Where(sender.ID(id)).
		Only(ctx)
}

func (r *senderRepo) Update(ctx context.Context, id uuid.UUID, name string, email *string) (*ent.Sender, error) {
	update := r.client.Sender.
		UpdateOneID(id).
		SetName(name)

	if email != nil {
		update = update.SetEmail(*email)
	} else {
		update = update.ClearEmail()
	}

	return update.Save(ctx)
}

func (r *senderRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Sender.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *senderRepo) List(ctx context.Context) ([]*ent.Sender, error) {
	return r.client.Sender.
		Query().
		All(ctx)
}

func (r *senderRepo) ListByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Sender, error) {
	return r.client.Sender.
		Query().
		Where(sender.CompanyID(companyID)).
		All(ctx)
}
