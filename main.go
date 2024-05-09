package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/milwad-dev/do-it/internal/handlers"
	"github.com/milwad-dev/do-it/internal/routers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Load environments
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Port
	port := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))

	// Connect to db and set it
	db := connectDB(err)

	defer db.Close()

	// Run tables
	runTables(db)

	// Set and get handler
	handler := handlers.NewDBHandler(db)

	// Router
	r := routers.GetRouter(handler)

	// Serve application
	fmt.Printf("Your application run on %v \n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal(err)
	}
}

// connectDB => connect to DB
func connectDB(err error) *sql.DB {
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_DATABASE")

	dns := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dns)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func runTables(db *sql.DB) {
	tables := make(map[string]string)

	tables["users"] = `CREATE TABLE IF NOT EXISTS users (
    	id INT NOT NULL AUTO_INCREMENT UNIQUE,
    	name VARCHAR(255) NOT NULL,
    	email VARCHAR(255) NULL UNIQUE,
    	phone VARCHAR(255) NULL UNIQUE,
    	password VARCHAR(255) NULL,
    	email_verified_at TIMESTAMP NULL DEFAULT NULL,
    	phone_verified_at TIMESTAMP NULL DEFAULT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	tables["labels"] = `CREATE TABLE IF NOT EXISTS labels (
    	id INT NOT NULL AUTO_INCREMENT UNIQUE,
    	title VARCHAR(255) NOT NULL,
    	user_id INT NOT NULL,
    	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`

	for _, table := range tables {
		db.Exec(table)
	}
}
