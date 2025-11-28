package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type UpdateHandler struct {
	documentService service.DocumentService
}

func NewUpdateHandler(documentService service.DocumentService) *UpdateHandler {
	return &UpdateHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Обновление метаданных документа
// @Description  Обновляет название, папку или отправителя документа
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID документа" format:"uuid"
// @Param        request body UpdateRequest true "Данные для обновления"
// @Success      200 {object} DocumentResponse "Документ успешно обновлен"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Документ не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/{id} [put]
func (h *UpdateHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	documentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	var req UpdateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	input := service.DocumentUpdateInput{
		Name:     req.Name,
		FolderID: req.FolderID,
		SenderID: req.SenderID,
	}

	document, err := h.documentService.Update(c.Context(), documentID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(DocumentResponse{
		ID:              document.ID,
		CompanyID:       document.CompanyID,
		FolderID:        document.FolderID,
		SenderID:        document.SenderID,
		FilePath:        document.FilePath,
		PreviewFilePath: document.PreviewFilePath,
		FileSize:        document.FileSize,
		MimeType:        document.MimeType,
		Checksum:        document.Checksum,
		CreatedAt:       document.CreatedAt,
	})
}
