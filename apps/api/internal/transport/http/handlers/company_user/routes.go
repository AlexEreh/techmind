package company_user

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes регистрирует маршруты для работы с пользователями компании
func RegisterRoutes(router fiber.Router, companyUserService service.CompanyUserService) {
	getMyCompaniesHandler := NewGetMyCompaniesHandler(companyUserService)

	router.Get("/my", getMyCompaniesHandler.Handle)
}
