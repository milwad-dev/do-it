package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
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

	// Router
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("milwad"))
	})

	// Serve application
	fmt.Printf("Your application run on %v", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
