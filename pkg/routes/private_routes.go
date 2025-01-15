package routes

import (
	"github.com/Szerl0k/go-office-desks/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Patch("/desks/:id", controllers.BookDesk)  // Book desks
	route.Post("/desks/add", controllers.CreateDesk) // Add desks
}
