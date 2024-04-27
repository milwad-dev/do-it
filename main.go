package main

import (
	"database/sql"
	"do-it/internal/handlers"
	"do-it/internal/routers"
	"fmt"
	"github.com/joho/godotenv"
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
	db, err := sql.Open("mysql", "")
	if err != nil {
		log.Fatal(err)
	}

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