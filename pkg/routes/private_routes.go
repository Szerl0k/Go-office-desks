package routes

import (
	"github.com/Szerl0k/go-office-desks/app/controllers"
	"github.com/Szerl0k/go-office-desks/pkg/middleware"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {

	a.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(middleware.GetJWTSecret())},
		ErrorHandler: middleware.JWTError,
	}))

	userRoutes(a)
	adminRoutes(a)

}

func userRoutes(a *fiber.App) {

	route := a.Group("/api/v1")

	route.Post("/user", controllers.GetUser)
	route.Get("/desks/", controllers.FetchAllDesks)
	route.Patch("/desks/:id", controllers.BookDesk)

}

func adminRoutes(a *fiber.App) {

	admin := a.Group("api/v1/admin", middleware.AdminMiddleware)

	admin.Post("/register", controllers.CreateUser)
	admin.Get("/user/:email", controllers.GetUserAdmin)
	admin.Get("/user", controllers.GetAllUsers)

	admin.Delete("/user/delete/:email", controllers.DeleteUser)
	admin.Post("/desks/add", controllers.CreateDesk)
	admin.Delete("/desks/delete/:id", controllers.DeleteDesk)

}
