package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

// InitRoutes initializes routes
func InitRoutes() http.Handler {
	r := chi.NewRouter()
	r.Get("/", status)

	return r
}

func status(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome home"))
}
