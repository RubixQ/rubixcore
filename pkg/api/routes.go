package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/rubixq/rubixcore/pkg/db"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"

	"github.com/go-chi/chi"
)

// InitRoutes initializes routes
func InitRoutes(l *zap.Logger, s *mgo.Session) http.Handler {
	r := chi.NewRouter()
	r.Post("/queues", createQueue(l, s))
	r.Get("/queues", listQueues(l, s))

	return r
}

func createQueue(l *zap.Logger, s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queue := db.Queue{
			ID: bson.NewObjectId(),
		}

		err := json.NewDecoder(r.Body).Decode(&queue)
		if err != nil {
			l.Error("failed decoding request payload", zap.Any("error", err))
			return
		}

		session := s.Copy()
		defer session.Close()

		queue.Active = true
		queue.Title = strings.ToUpper(strings.Replace(queue.Name, " ", "", -1))

		err = session.DB("rubixcore").C("queues").Insert(queue)
		if err != nil {
			l.Error("failed inserting queue", zap.Any("error", err))
			return
		}

		Ok(w, queue)
	}
}

func listQueues(l *zap.Logger, s *mgo.Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var queues []db.Queue

		err := session.DB("rubixcore").C("queues").Find(nil).All(&queues)
		if err != nil {
			l.Error("failed fetching queues from db", zap.Any("error", err))
			InternalServerError(w)
			return
		}

		Ok(w, queues)

	}
}
