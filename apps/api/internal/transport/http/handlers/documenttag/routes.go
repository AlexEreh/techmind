package documenttag

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с тегами документов
func RegisterRoutes(router fiber.Router, documentTagService service.DocumentTagService) {
	getDocumentTagsHandler := NewGetDocumentTagsHandler(documentTagService)
	addTagHandler := NewAddTagHandler(documentTagService)
	removeTagHandler := NewRemoveTagHandler(documentTagService)
	createTagHandler := NewCreateTagHandler(documentTagService)
	deleteTagHandler := NewDeleteTagHandler(documentTagService)
	getTagsByCompanyHandler := NewGetTagsByCompanyHandler(documentTagService)
	getTagByIDHandler := NewGetTagByIDHandler(documentTagService)
	updateTagHandler := NewUpdateTagHandler(documentTagService)

	// Операции с тегами документов
	router.Get("/document/:document_id", getDocumentTagsHandler.Handle)
	router.Post("/add", addTagHandler.Handle)
	router.Post("/remove", removeTagHandler.Handle)

	// Управление тегами
	router.Post("/tags", createTagHandler.Handle)
	router.Get("/tags/:id", getTagByIDHandler.Handle)
	router.Put("/tags/:id", updateTagHandler.Handle)
	router.Delete("/tags/:id", deleteTagHandler.Handle)
	router.Get("/company/:company_id", getTagsByCompanyHandler.Handle)
}
