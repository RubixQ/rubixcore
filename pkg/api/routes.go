package api

import (
	"encoding/json"
	"net/http"

	"github.com/rubixq/rubixcore/pkg/db"
	"github.com/rubixq/rubixcore/pkg/db/repo"

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
		queue := new(db.Queue)

		err := json.NewDecoder(r.Body).Decode(queue)
		if err != nil {
			l.Error("failed decoding request payload", zap.Any("error", err))
			return
		}

		session := s.Copy()
		defer session.Close()

		repo := repo.NewQueueRepo(session)

		queue, err = repo.Create(queue)
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

		repo := repo.NewQueueRepo(session)

		queues, err := repo.FindAll()

		if err != nil {
			l.Error("failed fetching queues from db", zap.Any("error", err))
			InternalServerError(w)
			return
		}

		Ok(w, queues)

	}
}
