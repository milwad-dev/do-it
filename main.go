package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func main() {
	// Port
	port := fmt.Sprintf(":%v", os.Getenv("APP_PORT"))

	// Router
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("milwad"))
	})

	// Serve application
	fmt.Printf("Your application run on %v", port)
	err := http.ListenAndServe(port, router)
	if err != nil {
		log.Fatal(err)
	}
}
