package routers

import (
	"do-it/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func GetRouter(handler *handlers.DBHandler) *chi.Mux {
	router := chi.NewRouter()

	// Users
	router.Get("/users", handler.GetLatestUsers)

	return router
}
