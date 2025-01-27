package main

import (
	"github.com/Szerl0k/go-office-desks/pkg/middleware"
	"github.com/Szerl0k/go-office-desks/pkg/routes"
	_ "github.com/go-sql-driver/mysql"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	port string
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port = os.Getenv("PORT")

}

func main() {

	app := fiber.New()

	routes.PublicRoutes(app)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(middleware.GetJWTSecret())},
		ErrorHandler: middleware.JWTError,
	}))

	routes.PrivateRoutes(app)

	log.Fatal(app.Listen(":" + port))

}
