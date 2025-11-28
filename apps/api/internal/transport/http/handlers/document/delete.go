package document

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteHandler struct {
	documentService service.DocumentService
}

func NewDeleteHandler(documentService service.DocumentService) *DeleteHandler {
	return &DeleteHandler{
		documentService: documentService,
	}
}

// Handle godoc
// @Summary      Удаление документа
// @Description  Удаляет документ и его файлы из хранилища
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID документа" format:"uuid"
// @Success      204 "Документ успешно удален"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Документ не найден"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/documents/{id} [delete]
func (h *DeleteHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	documentID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid document id format",
		})
	}

	if err := h.documentService.Delete(c.Context(), documentID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
