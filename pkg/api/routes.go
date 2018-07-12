package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// InitRoutes initializes routes
func InitRoutes(s *mgo.Session, l *zap.Logger) http.Handler {
	r := chi.NewRouter()
	r.Post("/queues", createQueue(s, l))
	r.Get("/queues", listQueues(s, l))

	return r
}
