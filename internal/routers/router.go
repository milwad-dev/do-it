package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/internal/handlers"
)

func GetRouter(handler *handlers.DBHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		// Users
		r.Get("/users", handler.GetLatestUsers)

		// Labels
		r.Post("/labels", handler.StoreLabel)

		// Tasks
		r.Get("/tasks", handler.GetLatestTasks)
		r.Post("/tasks", handler.StoreTask)

		// Auth
		r.Post("/register", handler.RegisterAuth)
		r.Post("/login", handler.LoginAuth)
	})

	return router
}
