package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
)

// InitRoutes initializes routes
func InitRoutes(s *mgo.Session, l *zap.Logger) http.Handler {
	r := mux.NewRouter()
	r.Handle("/queues", createQueue(s, l)).Methods("POST")
	r.Handle("/queues", listQueues(s, l)).Methods("GET")
	r.HandleFunc("/ws/status", handleStatusCheck(l)).Methods("GET")
	r.HandleFunc("/ws/test", handleStatusTest(l)).Methods("GET")

	return r
}
