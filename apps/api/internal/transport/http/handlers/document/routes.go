package document

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с документами
func RegisterRoutes(router fiber.Router, documentService service.DocumentService) {
	uploadHandler := NewUploadHandler(documentService)
	getByIDHandler := NewGetByIDHandler(documentService)
	getByFolderHandler := NewGetByFolderHandler(documentService)
	getByCompanyHandler := NewGetByCompanyHandler(documentService)
	updateHandler := NewUpdateHandler(documentService)
	deleteHandler := NewDeleteHandler(documentService)
	getDownloadURLHandler := NewGetDownloadURLHandler(documentService)
	getPreviewURLHandler := NewGetPreviewURLHandler(documentService)
	searchHandler := NewSearchHandler(documentService)
	router.Post("/", uploadHandler.Handle)
	router.Get("/:id", getByIDHandler.Handle)
	router.Put("/:id", updateHandler.Handle)
	router.Delete("/:id", deleteHandler.Handle)
	router.Get("/:id/download", getDownloadURLHandler.Handle)
	router.Get("/:id/preview", getPreviewURLHandler.Handle)
	router.Get("/folder/:folder_id", getByFolderHandler.Handle)
	router.Get("/company/:company_id", getByCompanyHandler.Handle)
	router.Post("/search", searchHandler.Handle)
}
