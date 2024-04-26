package routers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("milwad"))
	})

	return router
}
