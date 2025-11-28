package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type DeleteHandler struct {
	folderService service.FolderService
}

func NewDeleteHandler(folderService service.FolderService) *DeleteHandler {
	return &DeleteHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Удаление папки
// @Description  Удаляет папку и все вложенные папки и документы
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID папки" format:"uuid"
// @Success      204 "Папка успешно удалена"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Папка не найдена"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders/{id} [delete]
func (h *DeleteHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	folderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid folder id format",
		})
	}

	if err := h.folderService.Delete(c.Context(), folderID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
