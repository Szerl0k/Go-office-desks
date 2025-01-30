package controllers

import (
	"database/sql"
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

func scanUserByEmail(db *sql.DB, dest *structs.User, email string) error {

	query := "SELECT Email, Name, Surname, Password, Role FROM " + tUser + " WHERE Email Like ?"

	rows, err := db.Query(query, email)

	if err != nil {
		return err
	}

	rows.Next()
	return rows.Scan(&dest.Email, &dest.Name, &dest.Surname, &dest.Password, &dest.Role)
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

	user := &structs.User{}

	if err = scanUserByEmail(db, user, data.Email); err != nil {
		log.Printf("error fetching user data")
		return errorResponse(c, errors.New("invalid username or password"))
	}

	if user.Email == "" || !checkPasswordHash(user.Password, data.Password) {
		log.Printf("Invalid username or password: %v", data.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "invalid username or password",
		})
	}

	token, err := middleware.GenerateJWT(user.Email, user.Role)

	if err != nil {
		return errorResponse(c, errors.New("failed to generate token"))
	}

	return c.JSON(fiber.Map{
		"token": token,
	})

}

func CreateUser(c *fiber.Ctx) error {

	data := &structs.User{}

	if err := c.BodyParser(data); err != nil {
		return errorResponse(c, err)
	}

	if !data.ValidEmail() {
		return errorResponse(c, errors.New("incorrect email format"))
	}

	if data.Role == "" {
		data.Role = "user"
	}

	hashedPassword, err := hashPassword(data.Password)

	if err != nil {
		return errorResponse(c, errors.New("failed to hash password"))
	}

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	query := "INSERT INTO " + tUser + " (Email, Name, Surname, Password, Role) VALUES (?,?,?,?,?)"

	_, err = db.Exec(query, &data.Email, &data.Name, &data.Surname, hashedPassword, &data.Role)

	if err != nil {
		log.Println(err)
		return errorResponse(c, errors.New("email is already in use"))
	}

	// No need to send this information in response
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

	userEmail := claims["user_email"].(string)

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	userDetail := &structs.User{}

	if err := scanUserByEmail(db, userDetail, userEmail); err != nil {
		log.Printf("error scanning user row: %v", err)
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

func GetUserAdmin(c *fiber.Ctx) error {

	email := c.Params("email")
	user := &structs.User{}

	db, err := database.MySQLConnection()
	if err != nil {
		return errorResponse(c, err)
	}

	if err = scanUserByEmail(db, user, email); err != nil {
		log.Printf("error scanning user row: %v", err)
		return errorResponse(c, errors.New("user does not exist"))
	}

	user.Password = ""

	return c.JSON(fiber.Map{
		"error": "false",
		"msg":   "user found",
		"user":  user,
	})

}

func GetAllUsers(c *fiber.Ctx) error {

	var users []structs.User

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	rows, err := db.Query("SELECT Email, Name, Surname, Role FROM User")

	if err != nil {
		log.Printf(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "true",
			"msg":   "no users were found",
			"count": 0,
			"users": nil,
		})
	}

	for rows.Next() {
		var email string
		var name string
		var surname string
		var role string

		if err := rows.Scan(&email, &name, &surname, &role); err != nil {
			return errorResponse(c, err)
		}

		users = append(users, structs.User{
			Email:    email,
			Name:     name,
			Surname:  surname,
			Password: "",
			Role:     role,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"msg":   "users found",
		"count": len(users),
		"users": users,
	})

}

func DeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	query := "DELETE FROM " + tUser + " WHERE Email = ?"

	result, err := db.Exec(query, email)

	if err != nil {
		log.Printf("error deleting user")
		return errorResponse(c, err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errorResponse(c, errors.New("user does not exist"))
	}

	return c.Status(200).JSON(fiber.Map{
		"error": "false",
		"msg":   "Successfully deleted user",
		"email": email,
	})

}
