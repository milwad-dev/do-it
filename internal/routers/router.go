package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/docs"
	"github.com/milwad-dev/do-it/internal/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
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

		// Swagger
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8000/api/debug/swagger"), // The URL pointing to API definition
		))
		r.Get("/debug/swagger", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(docs.SwaggerInfo.SwaggerTemplate))
		})
	})

	return router
}
