package company_user

import (
	"techmind/internal/service"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
)

type GetMyCompaniesHandler struct {
	companyUserService service.CompanyUserService
}

func NewGetMyCompaniesHandler(companyUserService service.CompanyUserService) *GetMyCompaniesHandler {
	return &GetMyCompaniesHandler{
		companyUserService: companyUserService,
	}
}

// Handle godoc
// @Summary      Получение списка компаний пользователя
// @Description  Возвращает список всех компаний, в которых состоит текущий пользователь, с информацией о ролях
// @Tags         companies
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} MyCompaniesResponse "Список компаний"
// @Failure      401 {object} handlers.ErrorResponse "Неавторизированный доступ"
// @Failure      500 {object} handlers.ErrorResponse "Внутренняя ошибка сервера"
// @Router       /private/companies/my [get]
func (h *GetMyCompaniesHandler) Handle(c fiber.Ctx) error {
	userID, err := handlers.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	companies, err := h.companyUserService.GetUserCompanies(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(handlers.ErrorResponse{
			Error: err.Error(),
		})
	}

	response := MyCompaniesResponse{
		Companies: make([]CompanyUserData, 0, len(companies)),
	}

	for _, cu := range companies {
		companyData := CompanyUserData{
			ID:        cu.ID,
			UserID:    cu.UserID,
			CompanyID: cu.CompanyID,
			Role:      cu.Role,
		}

		// Добавляем информацию о компании из edges
		if cu.Edges.Company != nil {
			companyData.Company = &CompanyData{
				ID:   cu.Edges.Company.ID,
				Name: cu.Edges.Company.Name,
			}
		}

		response.Companies = append(response.Companies, companyData)
	}

	return c.JSON(response)
}
