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
		r.Delete("/labels/{id}", handler.DeleteLabel)

		// Tasks
		r.Get("/tasks", handler.GetLatestTasks)
		r.Post("/tasks", handler.StoreTask)
		r.Patch("/tasks/{id}/mark-as-completed", handler.MarkTaskAsCompleted)
		r.Delete("/tasks/{id}", handler.DeleteTask)
	})

	// Auth
	router.Post("/api/register", handler.RegisterAuth)
	router.Post("/api/login", handler.LoginAuth)
	router.Post("/api/forgot-password", handler.ForgotPasswordAuth)
	router.Post("/api/forgot-password-verify", handler.ForgotPasswordVerifyAuth)
	router.Post("/api/reset-password", handler.ResetPasswordAuth)

	// Swagger
	router.Get("/api/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/api/debug/swagger"), // The URL pointing to API definition
	))
	router.Get("/api/debug/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(docs.SwaggerInfo.SwaggerTemplate))
	})

	return router
}
