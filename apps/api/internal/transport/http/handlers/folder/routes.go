package folder

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с папками
func RegisterRoutes(router fiber.Router, folderService service.FolderService) {
	createHandler := NewCreateHandler(folderService)
	deleteHandler := NewDeleteHandler(folderService)
	renameHandler := NewRenameHandler(folderService)
	getByIDHandler := NewGetByIDHandler(folderService)
	getByCompanyHandler := NewGetByCompanyHandler(folderService)
	getByParentHandler := NewGetByParentHandler(folderService)

	router.Post("/", createHandler.Handle)
	router.Get("/:id", getByIDHandler.Handle)
	router.Delete("/:id", deleteHandler.Handle)
	router.Put("/:id/rename", renameHandler.Handle)
	router.Get("/company/:company_id", getByCompanyHandler.Handle)
	router.Post("/by-parent", getByParentHandler.Handle)
}
