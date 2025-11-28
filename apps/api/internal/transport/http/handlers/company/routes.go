package company

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с компаниями
func RegisterRoutes(router fiber.Router, companyService service.CompanyService) {
	createCompanyHandler := NewCreateCompanyHandler(companyService)

	router.Post("/", createCompanyHandler.Handle)
}
