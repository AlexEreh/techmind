package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type GetByParentHandler struct {
	folderService service.FolderService
}

func NewGetByParentHandler(folderService service.FolderService) *GetByParentHandler {
	return &GetByParentHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Получение вложенных папок
// @Description  Возвращает список папок по родительской папке. Если parent_id не указан, возвращает корневые папки
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body GetByParentRequest true "Параметры запроса"
// @Success      200 {object} FoldersListResponse "Список папок"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат запроса"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders/by-parent [post]
func (h *GetByParentHandler) Handle(c fiber.Ctx) error {
	var req GetByParentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid request format",
		})
	}

	folders, err := h.folderService.GetByParent(c.Context(), req.CompanyID, req.ParentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := FoldersListResponse{
		Folders: make([]FolderResponse, 0, len(folders)),
		Total:   len(folders),
	}

	for _, folder := range folders {
		response.Folders = append(response.Folders, FolderResponse{
			ID:             folder.ID,
			CompanyID:      folder.CompanyID,
			ParentFolderID: folder.ParentFolderID,
			Name:           folder.Name,
			Size:           folder.Size,
			Count:          folder.Count,
		})
	}

	return c.JSON(response)
}
