package routes

import (
	"github.com/Szerl0k/go-office-desks/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/login", controllers.Login)

}
