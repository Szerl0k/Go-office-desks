package routes

import (
	"github.com/Szerl0k/go-office-desks/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Patch("/desks/:id", controllers.BookDesk)         // Book desk
	route.Post("/desks/add", controllers.CreateDesk)        // Create desk
	route.Post("/desks/remove/:id", controllers.DeleteDesk) // Remove desks

	route.Post("/user", controllers.GetUser)
	route.Post("/register", controllers.CreateUser)

}
