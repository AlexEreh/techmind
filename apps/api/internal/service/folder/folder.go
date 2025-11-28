package folder

import (
	"context"
	"fmt"
	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

type folderService struct {
	folderRepo repo.FolderRepository
}

func NewService(folderRepo repo.FolderRepository) service.FolderService {
	return &folderService{
		folderRepo: folderRepo,
	}
}

func (s *folderService) Create(ctx context.Context, companyID uuid.UUID, name string, parentID *uuid.UUID) (*ent.Folder, error) {
	// Если указан parentID, проверяем что родительская папка существует
	if parentID != nil {
		parent, err := s.folderRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("parent folder not found: %w", err)
		}

		// Проверяем что родительская папка принадлежит той же компании
		if parent.CompanyID != companyID {
			return nil, fmt.Errorf("parent folder belongs to different company")
		}
	}

	// Создаем папку
	folder, err := s.folderRepo.Create(ctx, companyID, parentID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create folder: %w", err)
	}

	return folder, nil
}

func (s *folderService) Delete(ctx context.Context, folderID uuid.UUID) error {
	// Проверяем что папка существует
	_, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return fmt.Errorf("folder not found: %w", err)
	}

	// Удаляем папку (каскадное удаление вложенных папок и документов должно быть на уровне БД)
	if err := s.folderRepo.Delete(ctx, folderID); err != nil {
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}

func (s *folderService) Rename(ctx context.Context, folderID uuid.UUID, newName string) (*ent.Folder, error) {
	// Получаем текущую папку
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	// Обновляем имя папки, сохраняем текущие size и count
	updatedFolder, err := s.folderRepo.Update(ctx, folderID, newName, folder.Size, folder.Count)
	if err != nil {
		return nil, fmt.Errorf("failed to rename folder: %w", err)
	}

	return updatedFolder, nil
}

func (s *folderService) GetByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Folder, error) {
	folders, err := s.folderRepo.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders by company: %w", err)
	}

	return folders, nil
}

func (s *folderService) GetByParent(ctx context.Context, companyID uuid.UUID, parentID *uuid.UUID) ([]*ent.Folder, error) {
	// Если parentID nil, получаем корневые папки компании
	if parentID == nil {
		folders, err := s.folderRepo.ListByCompany(ctx, companyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get root folders: %w", err)
		}

		// Фильтруем только корневые папки (без родителя)
		var rootFolders []*ent.Folder
		for _, folder := range folders {
			if folder.ParentFolderID == nil {
				rootFolders = append(rootFolders, folder)
			}
		}

		return rootFolders, nil
	}

	// Проверяем что родительская папка существует и принадлежит компании
	parent, err := s.folderRepo.GetByID(ctx, *parentID)
	if err != nil {
		return nil, fmt.Errorf("parent folder not found: %w", err)
	}

	if parent.CompanyID != companyID {
		return nil, fmt.Errorf("parent folder belongs to different company")
	}

	// Получаем дочерние папки
	folders, err := s.folderRepo.ListByParent(ctx, *parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get folders by parent: %w", err)
	}

	return folders, nil
}

func (s *folderService) GetByID(ctx context.Context, folderID uuid.UUID) (*ent.Folder, error) {
	folder, err := s.folderRepo.GetByID(ctx, folderID)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}

	return folder, nil
}
