package routers

import (
	"do-it/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetRouter(handler *handlers.DBHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("milwad"))
	})

	return router
}
