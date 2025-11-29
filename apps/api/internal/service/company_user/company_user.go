package company_user

import (
	"context"
	"fmt"

	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

type companyUserService struct {
	repo repo.CompanyUserRepository
}

func NewService(repo repo.CompanyUserRepository) service.CompanyUserService {
	return &companyUserService{
		repo: repo,
	}
}

func (s *companyUserService) GetUserRole(ctx context.Context, userID, companyID uuid.UUID) (int, error) {
	role, err := s.repo.GetUserRole(ctx, userID, companyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user role: %w", err)
	}
	return role, nil
}

func (s *companyUserService) GetUserCompanies(ctx context.Context, userID uuid.UUID) ([]*ent.CompanyUser, error) {
	companies, err := s.repo.ListByUserWithCompany(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user companies: %w", err)
	}
	return companies, nil
}

func (s *companyUserService) GetCompanyUsers(ctx context.Context, companyID uuid.UUID) ([]*ent.CompanyUser, error) {
	users, err := s.repo.ListByCompanyWithUser(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company users: %w", err)
	}
	return users, nil
}
