package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type SearchHandler struct {
	documentService service.DocumentService
}

func NewSearchHandler(documentService service.DocumentService) *SearchHandler {
	return &SearchHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Поиск документов
// @Description  Ищет документы по различным критериям: текстовый запрос, папка, теги
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body SearchRequest true "Параметры поиска"
// @Success      200 {object} DocumentsListResponse "Результаты поиска"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/search [post]
func (h *SearchHandler) Handle(c fiber.Ctx) error {
	var req SearchRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	docsWithTags, err := h.documentService.Search(c.Context(), req.CompanyID, req.Query, req.FolderID, req.TagIDs)
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
			FilePath:        docWithTags.Document.FilePath,
			PreviewFilePath: docWithTags.Document.PreviewFilePath,
			FileSize:        docWithTags.Document.FileSize,
			MimeType:        docWithTags.Document.MimeType,
			Checksum:        docWithTags.Document.Checksum,
			CreatedAt:       docWithTags.Document.CreatedAt,
			Tags:            tags,
			PreviewURL:      docWithTags.PreviewURL,
		})
	}

	return c.JSON(response)
}
