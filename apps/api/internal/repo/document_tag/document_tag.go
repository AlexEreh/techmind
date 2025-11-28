package document_tag

import (
	"context"
	"techmind/internal/repo"
	"techmind/schema/ent"
	"techmind/schema/ent/documenttag"

	"github.com/google/uuid"
)

type documentTagRepo struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) repo.DocumentTagRepository {
	return &documentTagRepo{client: client}
}

func (r *documentTagRepo) Create(ctx context.Context, documentID, tagID uuid.UUID) (*ent.DocumentTag, error) {
	return r.client.DocumentTag.
		Create().
		SetDocumentID(documentID).
		SetTagID(tagID).
		Save(ctx)
}

func (r *documentTagRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.DocumentTag, error) {
	return r.client.DocumentTag.
		Query().
		Where(documenttag.ID(id)).
		Only(ctx)
}

func (r *documentTagRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.DocumentTag.
		DeleteOneID(id).
		Exec(ctx)
}

func (r *documentTagRepo) List(ctx context.Context) ([]*ent.DocumentTag, error) {
	return r.client.DocumentTag.
		Query().
		All(ctx)
}

func (r *documentTagRepo) ListByDocument(ctx context.Context, documentID uuid.UUID) ([]*ent.DocumentTag, error) {
	return r.client.DocumentTag.
		Query().
		Where(documenttag.DocumentID(documentID)).
		All(ctx)
}

func (r *documentTagRepo) ListByTag(ctx context.Context, tagID uuid.UUID) ([]*ent.DocumentTag, error) {
	return r.client.DocumentTag.
		Query().
		Where(documenttag.TagID(tagID)).
		All(ctx)
}

func (r *documentTagRepo) DeleteByDocumentAndTag(ctx context.Context, documentID, tagID uuid.UUID) error {
	_, err := r.client.DocumentTag.
		Delete().
		Where(
			documenttag.DocumentID(documentID),
			documenttag.TagID(tagID),
		).
		Exec(ctx)
	return err
}
