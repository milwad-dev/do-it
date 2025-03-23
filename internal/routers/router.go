package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/milwad-dev/do-it/docs"
	"github.com/milwad-dev/do-it/internal/handlers"
	"github.com/milwad-dev/do-it/internal/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func GetRouter(handler *handlers.DBHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.Authenticate)

		// Users
		r.Get("/users", handler.GetLatestUsers)

		// Labels
		r.Get("/labels", handler.GetLatestLabels)
		r.Post("/labels", handler.StoreLabel)

		// Tasks
		r.Get("/tasks", handler.GetLatestTasks)
		r.Post("/tasks", handler.StoreTask)
		r.Patch("/tasks/{task}/mark-as-completed", handler.MarkTaskAsCompleted)

		// Swagger
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8000/api/debug/swagger"), // The URL pointing to API definition
		))
		r.Get("/debug/swagger", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(docs.SwaggerInfo.SwaggerTemplate))
		})
	})

	router.Route("/api", func(r chi.Router) {
		// Auth
		r.Post("/register", handler.RegisterAuth)
		r.Post("/login", handler.LoginAuth)
	})

	return router
}
