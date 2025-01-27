package controllers

import (
	"errors"
	"github.com/Szerl0k/go-office-desks/app/structs"
	"github.com/Szerl0k/go-office-desks/pkg/middleware"
	"github.com/Szerl0k/go-office-desks/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {

	data := &structs.User{}

	if err := c.BodyParser(data); err != nil {
		return errorResponse(c, err)
	}

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	query := "SELECT Email, Name, Surname, Password  FROM " + tUser + " WHERE Email Like ?"

	rows, err := db.Query(query, data.Email)

	if err != nil {
		log.Printf("error fetching user data")
		return errorResponse(c, errors.New("invalid username or password"))
	}

	user := structs.User{}

	rows.Next()

	if err := rows.Scan(&user.Email, &user.Name, &user.Surname, &user.Password); err != nil {
		log.Printf("error scanning user row")
		return errorResponse(c, errors.New("invalid username or password 2"))
	}

	if user.Email == "" { // || !checkPasswordHash(user.Password, data.Password)
		log.Printf("Invalid username or password: %v", data.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid username or password",
		})
	}

	token, err := middleware.GenerateJWT(user.Email)

	if err != nil {
		return errorResponse(c, errors.New("failed to generate token"))
	}

	return c.JSON(fiber.Map{
		"token": token,
	})

}

func CreateUser(c *fiber.Ctx) error {

	data := structs.User{}

	if err := c.BodyParser(data); err != nil {
		return errorResponse(c, err)
	}

	hashedPasword, err := hashPassword(data.Password)

	if err != nil {
		return errorResponse(c, errors.New("failed to hash password"))
	}

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	query := "INSERT INTO " + tUser + " (Email, Name, Surname, Password) VALUES (?,?,?,?)"

	_, err = db.Exec(query, &data.Email, &data.Name, &data.Surname, hashedPasword)

	if err != nil {
		return errorResponse(c, errors.New("user can't be created"))
	}

	data.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "user created",
		"user":  data,
	})

}

func GetUser(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userEmail := claims["user_email"].(float64)

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	query := "SELECT Email, Name, Surname, password  FROM " + tUser + " WHERE Email Like ?"

	rows, err := db.Query(query, userEmail)

	if err != nil {
		log.Printf("error fetching user data")
		return errorResponse(c, errors.New("invalid username or password"))
	}

	userDetail := structs.User{}

	rows.Next()

	if err := rows.Scan(&userDetail.Email, &userDetail.Name, &userDetail.Surname, &userDetail.Password); err != nil {
		log.Printf("error scanning user row")
		return errorResponse(c, errors.New("invalid username or password 2"))
	}

	if userDetail.Email == "" {
		return errorResponse(c, errors.New("user not found"))
	}

	return c.JSON(fiber.Map{
		"error": false,
		"msg":   "user found",
		"user":  userDetail,
	})

}
