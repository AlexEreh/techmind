package folder

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByCompanyHandler struct {
	folderService service.FolderService
}

func NewGetByCompanyHandler(folderService service.FolderService) *GetByCompanyHandler {
	return &GetByCompanyHandler{
		folderService: folderService,
	}
}

// Handle godoc
// @Summary      Получение всех папок компании
// @Description  Возвращает список всех папок компании без учета иерархии
// @Tags         folders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        company_id path string true "ID компании" format:"uuid"
// @Success      200 {object} FoldersListResponse "Список папок"
// @Failure      400 {object} handlers.ErrorResponse "Неверный формат ID"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/folders/company/{company_id} [get]
func (h *GetByCompanyHandler) Handle(c fiber.Ctx) error {
	companyIDParam := c.Params("company_id")
	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid company id format",
		})
	}

	folders, err := h.folderService.GetByCompany(c.Context(), companyID)
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
