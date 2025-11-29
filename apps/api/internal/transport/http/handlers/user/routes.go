package user

import (
	"techmind/internal/repo"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, userRepo repo.UserRepository) {
	getByIDHandler := NewGetByIDHandler(userRepo)

	router.Get("/:id", getByIDHandler.Handle)
}
