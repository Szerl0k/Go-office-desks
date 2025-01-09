package main

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Desk struct {
	ID       int    `json:"id"`
	Occupied bool   `json:"occupied"`
	Body     string `json:"body"`
}

func main() {

	loadEnv()

	port := os.Getenv("PORT")

	db := connectToDb()
	defer db.Close()

	app := fiber.New()

	// GET DESKS
	app.Get("/api/desks/", func(c *fiber.Ctx) error {
		desks := fetchAllDesks(db)

		return c.Status(200).JSON(desks)
	})

	// UPDATE DESK
	app.Patch("/api/desks/:id", func(c *fiber.Ctx) error {

		id := c.Params("id")

		desk, err := fetchDesk(id, db)

		if err != nil {
			log.Printf("Attempt to update desk %v, but it is already occupied or it doesn't exist", id)
			return c.Status(409).JSON(fiber.Map{"error": "desk already occupied"})
		}
		return c.Status(200).JSON(desk)

	})

	// CREATE DESK
	/*	app.Post("/api/desks/", func(c *fiber.Ctx) error {
		desk := &Desk{}

		if err := c.BodyParser(desk); err != nil {
			return err
		}

		if desk.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body required"})
		}

		desk.ID = len(desks) + 1
		desk.Occupied = false
		desks = append(desks, *desk)

		return c.Status(201).JSON(desk)
	})*/

	log.Fatal(app.Listen(":" + port))

}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func connectToDb() *sql.DB {
	MysqlUri := os.Getenv("MYSQL_URI")

	db, err := sql.Open("mysql", MysqlUri)

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		log.Fatalf("Error verifying connection to the database: %v", err)
	}

	return db
}

func fetchAllDesks(db *sql.DB) []Desk {

	var desks []Desk

	rows, err := db.Query("SELECT * FROM Desk")

	if err != nil {
		log.Fatalf("Error querying the database: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var isOccupied bool
		var body string

		if err := rows.Scan(&id, &isOccupied, &body); err != nil {
			log.Fatalf("Error scanning row %v", err)
		}
		desks = append(desks, Desk{id, isOccupied, body})

	}

	return desks

}

func fetchDesk(id string, db *sql.DB) (Desk, error) {

	desk := Desk{}

	query := "UPDATE Desk SET occupied = ? WHERE id = ?"

	result, err := db.Exec(query, true, id)

	if err != nil {
		log.Fatalf("Error reserving desk: %v", err)
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return desk, errors.New("desk is already occupied or does not exist")
	}

	query = "SELECT * FROM Desk WHERE Id = ?"

	rows, err := db.Query(query, id)

	if err != nil {
		log.Fatalf("Error reserving desk: %v", err)
	}

	rows.Next()

	if err := rows.Scan(&desk.ID, &desk.Occupied, &desk.Body); err != nil {
		log.Fatalf("Error scanning row %v", err)
	}

	return desk, nil

}
