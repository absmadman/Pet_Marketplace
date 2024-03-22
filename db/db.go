package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InsertUser() {

}

func NewDBConnection() *sql.DB {
	pgPass := os.Getenv("POSTGRES_PASSWORD")
	pgUser := os.Getenv("POSTGRES_USER")
	pgDb := os.Getenv("POSTGRES_DB")
	connStr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", pgUser, pgPass, pgDb)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
