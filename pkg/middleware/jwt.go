package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var jwtSecret = os.Getenv("jwtSecret")

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"user_email": email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // valid 24hrs
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func GetJWTSecret() string {
	return jwtSecret
}

func JWTError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "Missing or malformed token",
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": true,
		"msg":   "Invalid or expired token",
	})
}
