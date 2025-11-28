package company

import (
	"fmt"
	"techmind/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type CreateCompanyHandler struct {
	companyService service.CompanyService
}

func NewCreateCompanyHandler(companyService service.CompanyService) *CreateCompanyHandler {
	return &CreateCompanyHandler{
		companyService: companyService,
	}
}

type CreateCompanyRequest struct {
	Name string `json:"name" validate:"required,min=1"`
}

func (h *CreateCompanyHandler) Handle(c fiber.Ctx) error {
	ctx := c.Context()

	fmt.Println("=== CreateCompany Handler ===")
	fmt.Printf("Authorization header: %s\n", c.Get("Authorization"))

	// Получаем ID пользователя из контекста (установлено middleware)
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		fmt.Println("ERROR: user_id not found in context or wrong type")
		fmt.Printf("Context locals: %+v\n", c.Locals("user_id"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	fmt.Printf("UserID from context: %s\n", userID)

	var req CreateCompanyRequest
	if err := c.Bind().JSON(&req); err != nil {
		fmt.Printf("ERROR: Failed to bind JSON: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	fmt.Printf("Request: %+v\n", req)

	if req.Name == "" {
		fmt.Println("ERROR: Company name is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "company name is required",
		})
	}

	company, err := h.companyService.Create(ctx, req.Name, userID)
	if err != nil {
		fmt.Printf("ERROR: Failed to create company: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create company",
		})
	}

	fmt.Printf("SUCCESS: Company created: %+v\n", company)
	return c.Status(fiber.StatusCreated).JSON(company)
}
