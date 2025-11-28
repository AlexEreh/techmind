package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type CreateHandler struct {
	folderService service.FolderService
}

func NewCreateHandler(folderService service.FolderService) *CreateHandler {
	return &CreateHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Создание папки
// @Description  Создает новую папку в компании, может быть вложенной
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateRequest true "Данные для создания папки"
// @Success      201 {object} FolderResponse "Папка успешно создана"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      404 {object} handlers.ErrorResponse "Родительская папка не найдена"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders [post]
func (h *CreateHandler) Handle(c fiber.Ctx) error {
	var req CreateRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	folder, err := h.folderService.Create(c.Context(), req.CompanyID, req.Name, req.ParentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(FolderResponse{
		ID:             folder.ID,
		CompanyID:      folder.CompanyID,
		ParentFolderID: folder.ParentFolderID,
		Name:           folder.Name,
		Size:           folder.Size,
		Count:          folder.Count,
	})
}
