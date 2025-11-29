package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByFolderHandler struct {
	documentService service.DocumentService
}

func NewGetByFolderHandler(documentService service.DocumentService) *GetByFolderHandler {
	return &GetByFolderHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Получение документов папки
// @Description  Возвращает список всех документов в указанной папке
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        folder_id path string true "ID папки" format:"uuid"
// @Success      200 {object} DocumentsListResponse "Список документов"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/folder/{folder_id} [get]
func (h *GetByFolderHandler) Handle(c fiber.Ctx) error {
	folderIDParam := c.Params("folder_id")
	folderID, err := uuid.Parse(folderIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid folder id format",
		})
	}

	docsWithTags, err := h.documentService.GetByFolder(c.Context(), folderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := DocumentsListResponse{
		Documents: make([]DocumentResponse, 0, len(docsWithTags)),
		Total:     len(docsWithTags),
	}

	for _, docWithTags := range docsWithTags {
		tags := make([]TagData, 0, len(docWithTags.Tags))
		for _, tag := range docWithTags.Tags {
			tags = append(tags, TagData{
				ID:        tag.ID,
				CompanyID: tag.CompanyID,
				Name:      tag.Name,
			})
		}

		response.Documents = append(response.Documents, DocumentResponse{
			ID:              docWithTags.Document.ID,
			CompanyID:       docWithTags.Document.CompanyID,
			FolderID:        docWithTags.Document.FolderID,
			SenderID:        docWithTags.Document.SenderID,
			Name:            docWithTags.Document.Name,
			FilePath:        docWithTags.Document.FilePath,
			PreviewFilePath: docWithTags.Document.PreviewFilePath,
			FileSize:        docWithTags.Document.FileSize,
			MimeType:        docWithTags.Document.MimeType,
			Checksum:        docWithTags.Document.Checksum,
			CreatedBy:       docWithTags.Document.CreatedBy,
			UpdatedBy:       docWithTags.Document.UpdatedBy,
			CreatedAt:       docWithTags.Document.CreatedAt,
			UpdatedAt:       docWithTags.Document.UpdatedAt,
			Tags:            tags,
			PreviewURL:      docWithTags.PreviewURL,
			DownloadURL:     docWithTags.DownloadURL,
		})
	}

	return c.JSON(response)
}
