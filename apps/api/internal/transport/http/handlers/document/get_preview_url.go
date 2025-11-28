package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetPreviewURLHandler struct {
	documentService service.DocumentService
}

func NewGetPreviewURLHandler(documentService service.DocumentService) *GetPreviewURLHandler {
	return &GetPreviewURLHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Получение ссылки на preview
// @Description  Возвращает временную presigned URL для доступа к preview документа
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID документа" format:"uuid"
// @Success      200 {object} URLResponse "Ссылка на preview"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Preview не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/{id}/preview [get]
func (h *GetPreviewURLHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	documentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	url, err := h.documentService.GetPreviewURL(c.Context(), documentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(URLResponse{
		URL:       url,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	})
}
