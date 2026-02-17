package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func ConnectDb(username string, password string, host string, port string, dbname string) (*sql.DB, error) {

	// Build connection string with Sprintf
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		username, password, host, port, dbname,
	)

	// Open a database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	defer db.Close() // Always close when done

	// Verify the connection is alive
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return db, nil
}
