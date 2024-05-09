package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/internal/handlers"
)

func GetRouter(handler *handlers.DBHandler) *chi.Mux {
	router := chi.NewRouter()

	// Users
	router.Get("/users", handler.GetLatestUsers)

	// Labels
	router.Post("/labels", handler.StoreLabel)

	// TODOS
	//router.Get("/todos", handler.GetLatestTodos)

	return router
}
