package sender

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с контрагентами
func RegisterRoutes(router fiber.Router, senderService service.SenderService) {
	createHandler := NewCreateHandler(senderService)
	getByIDHandler := NewGetByIDHandler(senderService)
	updateHandler := NewUpdateHandler(senderService)
	deleteHandler := NewDeleteHandler(senderService)
	getByCompanyHandler := NewGetByCompanyHandler(senderService)

	// CRUD операции с контрагентами
	router.Post("/", createHandler.Handle)
	router.Get("/company/:company_id", getByCompanyHandler.Handle)
	router.Get("/:id", getByIDHandler.Handle)
	router.Put("/:id", updateHandler.Handle)
	router.Delete("/:id", deleteHandler.Handle)
}
