package user

import (
	"techmind/internal/repo"
	"techmind/internal/transport/http/handlers"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type GetByIDHandler struct {
	userRepo repo.UserRepository
}

func NewGetByIDHandler(userRepo repo.UserRepository) *GetByIDHandler {
	return &GetByIDHandler{
		userRepo: userRepo,
	}
}

func (h *GetByIDHandler) Handle(c fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(handlers.ErrorResponse{
			Error: "invalid user id format",
		})
	}

	user, err := h.userRepo.GetByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(handlers.ErrorResponse{
			Error: "user not found",
		})
	}

	return c.JSON(UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}
