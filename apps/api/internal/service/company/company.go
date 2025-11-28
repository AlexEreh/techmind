package company

import (
	"context"

	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo     repo.CompanyRepository
	companyUserRepo repo.CompanyUserRepository
}

func NewService(companyRepo repo.CompanyRepository, companyUserRepo repo.CompanyUserRepository) service.CompanyService {
	return &CompanyService{
		companyRepo:     companyRepo,
		companyUserRepo: companyUserRepo,
	}
}

// Create создает новую компанию и автоматически добавляет создателя как администратора
func (s *CompanyService) Create(ctx context.Context, name string, userID uuid.UUID) (*ent.Company, error) {
	// Создаем компанию
	company, err := s.companyRepo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	// Добавляем создателя как администратора (роль 2)
	_, err = s.companyUserRepo.Create(ctx, userID, company.ID, 2)
	if err != nil {
		// Если не удалось добавить пользователя, удаляем компанию
		_ = s.companyRepo.Delete(ctx, company.ID)
		return nil, err
	}

	return company, nil
}
