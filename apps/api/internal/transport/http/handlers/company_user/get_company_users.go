package company_user

import (
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetCompanyUsersHandler struct {
	companyUserService service.CompanyUserService
}

func NewGetCompanyUsersHandler(companyUserService service.CompanyUserService) *GetCompanyUsersHandler {
	return &GetCompanyUsersHandler{
		companyUserService: companyUserService,
	}
}

func (h *GetCompanyUsersHandler) Handle(c fiber.Ctx) error {
	// Получаем ID компании из параметров маршрута
	companyIDParam := c.Params("companyId")
	companyID, err := uuid.Parse(companyIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid company ID",
		})
	}

	// Получаем список пользователей компании
	companyUsers, err := h.companyUserService.GetCompanyUsers(c.Context(), companyID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get company users",
		})
	}

	// Преобразуем в DTO для фронтенда
	users := make([]CompanyUserWithDetailsDTO, 0, len(companyUsers))
	for _, cu := range companyUsers {
		// Проверяем, что пользователь загружен
		if cu.Edges.User == nil {
			continue
		}

		users = append(users, CompanyUserWithDetailsDTO{
			ID:       cu.Edges.User.ID,
			Username: cu.Edges.User.Name,
			Email:    cu.Edges.User.Email,
			Role:     cu.Role,
			AddedAt:  cu.AddedAt,
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
		"total": len(users),
	})
}
