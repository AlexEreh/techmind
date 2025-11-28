package documenttag

import (
	"context"
	"fmt"
	"techmind/internal/repo"
	"techmind/internal/service"
	"techmind/schema/ent"

	"github.com/google/uuid"
)

type documentTagService struct {
	documentTagRepo repo.DocumentTagRepository
	tagRepo         repo.TagRepository
	documentRepo    repo.DocumentRepository
}

func NewService(
	documentTagRepo repo.DocumentTagRepository,
	tagRepo repo.TagRepository,
	documentRepo repo.DocumentRepository,
) service.DocumentTagService {
	return &documentTagService{
		documentTagRepo: documentTagRepo,
		tagRepo:         tagRepo,
		documentRepo:    documentRepo,
	}
}

func (s *documentTagService) GetDocumentTags(ctx context.Context, documentID uuid.UUID) ([]*ent.Tag, error) {
	// Проверяем что документ существует
	_, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("document not found: %w", err)
	}

	// Получаем связи документ-тег
	docTags, err := s.documentTagRepo.ListByDocument(ctx, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document tags: %w", err)
	}

	// Получаем полные объекты тегов
	tags := make([]*ent.Tag, 0, len(docTags))
	for _, dt := range docTags {
		tag, err := s.tagRepo.GetByID(ctx, dt.TagID)
		if err != nil {
			continue // Пропускаем удаленные теги
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (s *documentTagService) AddTagToDocument(ctx context.Context, documentID, tagID uuid.UUID) error {
	// Проверяем что документ существует
	document, err := s.documentRepo.GetByID(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}

	// Проверяем что тег существует
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Проверяем что тег принадлежит той же компании
	if tag.CompanyID != document.CompanyID {
		return fmt.Errorf("tag belongs to different company")
	}

	// Проверяем что связь еще не существует
	existingTags, err := s.documentTagRepo.ListByDocument(ctx, documentID)
	if err == nil {
		for _, dt := range existingTags {
			if dt.TagID == tagID {
				return nil // Связь уже существует, не ошибка
			}
		}
	}

	// Создаем связь
	_, err = s.documentTagRepo.Create(ctx, documentID, tagID)
	if err != nil {
		return fmt.Errorf("failed to add tag to document: %w", err)
	}

	return nil
}

func (s *documentTagService) RemoveTagFromDocument(ctx context.Context, documentID, tagID uuid.UUID) error {
	// Удаляем связь
	err := s.documentTagRepo.DeleteByDocumentAndTag(ctx, documentID, tagID)
	if err != nil {
		return fmt.Errorf("failed to remove tag from document: %w", err)
	}

	return nil
}

func (s *documentTagService) CreateTag(ctx context.Context, companyID uuid.UUID, name string) (*ent.Tag, error) {
	// Проверяем что тег с таким именем не существует в компании
	existingTag, err := s.tagRepo.GetByName(ctx, companyID, name)
	if err == nil && existingTag != nil {
		return existingTag, nil // Возвращаем существующий тег
	}

	// Создаем новый тег
	tag, err := s.tagRepo.Create(ctx, companyID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return tag, nil
}

func (s *documentTagService) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	// Проверяем что тег существует
	_, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return fmt.Errorf("tag not found: %w", err)
	}

	// Удаляем тег (каскадно удалятся все связи с документами)
	if err := s.tagRepo.Delete(ctx, tagID); err != nil {
		return fmt.Errorf("failed to delete tag: %w", err)
	}

	return nil
}

func (s *documentTagService) GetTagsByCompany(ctx context.Context, companyID uuid.UUID) ([]*ent.Tag, error) {
	tags, err := s.tagRepo.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tags by company: %w", err)
	}

	return tags, nil
}

func (s *documentTagService) GetTagByID(ctx context.Context, tagID uuid.UUID) (*ent.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	return tag, nil
}

func (s *documentTagService) UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*ent.Tag, error) {
	// Проверяем что тег существует
	_, err := s.tagRepo.GetByID(ctx, tagID)
	if err != nil {
		return nil, fmt.Errorf("tag not found: %w", err)
	}

	// Обновляем название тега
	updatedTag, err := s.tagRepo.Update(ctx, tagID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return updatedTag, nil
}
