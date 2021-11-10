// This file sets up a database object for queries
package database

import (
    "fmt"
    "log"
    "os"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"strconv"
)

// SetupDB loads the data from the .env file and sets up the database object.
func SetupDB() *sql.DB {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get key .env variables
	host := os.Getenv("DB_HOST")
	port,_ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	// Generate string
	dbinfo := fmt.Sprintf("host=%s port=%d user=%s "+ "password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {
		log.Fatal(err)
	}
	
	return db
}