package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetDownloadURLHandler struct {
	documentService service.DocumentService
}

func NewGetDownloadURLHandler(documentService service.DocumentService) *GetDownloadURLHandler {
	return &GetDownloadURLHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Получение ссылки на скачивание
// @Description  Возвращает временную presigned URL для скачивания оригинала документа
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID документа" format:"uuid"
// @Success      200 {object} URLResponse "Ссылка для скачивания"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Документ не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/{id}/download [get]
func (h *GetDownloadURLHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	documentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	url, err := h.documentService.GetDownloadURL(c.Context(), documentID)
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
