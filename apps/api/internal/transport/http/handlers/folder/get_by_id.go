package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByIDHandler struct {
	folderService service.FolderService
}

func NewGetByIDHandler(folderService service.FolderService) *GetByIDHandler {
	return &GetByIDHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Получение папки по ID
// @Description  Возвращает информацию о папке по её идентификатору
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "ID папки" format:"uuid"
// @Success      200 {object} FolderResponse "Данные папки"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      404 {object} handlers.ErrorResponse "Папка не найдена"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders/{id} [get]
func (h *GetByIDHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	folderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid folder id format",
		})
	}

	folder, err := h.folderService.GetByID(c.Context(), folderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(handlers.ErrorResponse{
			Error: "folder not found",
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
