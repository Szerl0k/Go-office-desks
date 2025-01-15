package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

func MySQLConnection() (*sql.DB, error) {

	mysqlUri := os.Getenv("MYSQL_URI")

	db, err := sql.Open("mysql", mysqlUri)

	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error verifying connection to the database: %v", err)
	}

	return db, nil

}
