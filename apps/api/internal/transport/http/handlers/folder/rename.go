package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type RenameHandler struct {
	folderService service.FolderService
}

func NewRenameHandler(folderService service.FolderService) *RenameHandler {
	return &RenameHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Переименование папки
// @Description  Изменяет название папки
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID папки" format:"uuid"
// @Param        request body RenameRequest true "Новое название"
// @Success      200 {object} FolderResponse "Папка успешно переименована"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Папка не найдена"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders/{id}/rename [put]
func (h *RenameHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	folderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid folder id format",
		})
	}

	var req RenameRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	folder, err := h.folderService.Rename(c.Context(), folderID, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.JSON(FolderResponse{
		ID:             folder.ID,
		CompanyID:      folder.CompanyID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name,
		Size:           folder.Size,
		Count:          folder.Count,
	})
}
