package controllers

import (
	"errors"
	"github.com/Szerl0k/go-office-desks/app/structs"
	"github.com/Szerl0k/go-office-desks/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	_             = godotenv.Load(".env") // For some reason without this the os.Gatenv() will assign no values
	tDesks        = os.Getenv("tDesks")
	tUsers string = os.Getenv("tUsers")
)

func errorResponse(c *fiber.Ctx, err error) error {
	log.Printf(err.Error())
	return c.Status(fiber.StatusInternalServerError).JSON(
		fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
}

func FetchAllDesks(c *fiber.Ctx) error {
	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	var desks []structs.Desk

	rows, err := db.Query("SELECT * FROM " + tDesks + "")

	defer rows.Close()

	if err != nil {
		log.Printf(err.Error())

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "no desks were found",
			"count": 0,
			"desks": nil,
		})
	}

	for rows.Next() {
		var id int
		var floor int
		var isOccupied bool
		var body string

		if err := rows.Scan(&id, &floor, &isOccupied, &body); err != nil {
			return errorResponse(c, err)
		}

		desks = append(desks, structs.Desk{
			ID:       id,
			Floor:    floor,
			Occupied: isOccupied,
			Body:     body,
		})

	}

	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(desks),
		"desks": desks,
	})

}

func BookDesk(c *fiber.Ctx) error {

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	id := c.Params("id")
	desk := structs.Desk{}

	query := "UPDATE " + tDesks + " SET occupied = ? WHERE id = ?"

	result, err := db.Exec(query, true, id)

	if err != nil {
		log.Printf("error booking desk")
		return errorResponse(c, err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return errorResponse(c, errors.New("desk already occupied or does not exist"))
	}

	query = "SELECT * FROM " + tDesks + " WHERE Id = ?"

	rows, err := db.Query(query, id)

	if err != nil {
		log.Printf("error booking desk")
		return errorResponse(c, err)
	}

	rows.Next()

	if err := rows.Scan(&desk.ID, &desk.Floor, &desk.Occupied, &desk.Body); err != nil {
		log.Printf("error scanning row")
		return errorResponse(c, err)
	}

	return c.Status(409).JSON(fiber.Map{
		"error": false,
		"msg":   "successfully booked a desk",
		"desks": desk,
	})
}

func CreateDesk(c *fiber.Ctx) error {

	db, err := database.MySQLConnection()

	if err != nil {
		return errorResponse(c, err)
	}

	desk := &structs.Desk{}

	if err := c.BodyParser(desk); err != nil {
		return errorResponse(c, err)
	}

	if desk.Body == "" {
		return errorResponse(c, errors.New("desk floor and body is required"))
	}

	query := "INSERT INTO " + tDesks + " (Floor, Occupied, Body) VALUES (?,?,?)"

	result, err := db.Exec(query, &desk.Floor, false, &desk.Body)

	if err != nil {
		return errorResponse(c, errors.New("desk can't be created"))
	}

	id, _ := result.LastInsertId() // LAST_INSERT_ID() will never return a sql error

	desk.ID = int(id)

	return c.Status(200).JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"desk":  desk,
	})

}
