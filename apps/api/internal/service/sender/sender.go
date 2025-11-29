package sender

import (
	"context"
	"fmt"
	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

type senderService struct {
	repo repo.SenderRepository
}

func NewService(repo repo.SenderRepository) service.SenderService {
	return &senderService{
		repo: repo,
	}
}

func (s *senderService) Create(ctx context.Context, companyID uuid.UUID, name string, email *string) (*ent.Sender, error) {
	sender, err := s.repo.Create(ctx, companyID, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create sender: %w", err)
	}
	return sender, nil
}

func (s *senderService) GetByID(ctx context.Context, id uuid.UUID) (*ent.Sender, error) {
	sender, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}
	return sender, nil
}

func (s *senderService) Update(ctx context.Context, id uuid.UUID, name string, email *string) (*ent.Sender, error) {
	sender, err := s.repo.Update(ctx, id, name, email)
	if err != nil {
		return nil, fmt.Errorf("failed to update sender: %w", err)
	}
	return sender, nil
}

func (s *senderService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete sender: %w", err)
	}
	return nil
}

func (s *senderService) GetByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Sender, error) {
	senders, err := s.repo.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get senders by company: %w", err)
	}
	return senders, nil
}
