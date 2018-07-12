package api

import (
	"encoding/json"
	"net/http"

	"github.com/rubixq/rubixcore/pkg/db"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

func createQueue(s *mgo.Session, l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queue := new(db.Queue)

		err := json.NewDecoder(r.Body).Decode(queue)
		if err != nil {
			l.Error("failed decoding request payload", zap.Any("error", err))
			return
		}

		session := s.Copy()
		defer session.Close()

		repo := db.NewQueueRepo(session)

		queue, err = repo.Create(queue)
		if err != nil {
			l.Error("failed inserting queue", zap.Any("error", err))
			return
		}

		Ok(w, queue)
	}
}

func listQueues(s *mgo.Session, l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		repo := db.NewQueueRepo(session)

		queues, err := repo.FindAll()

		if err != nil {
			l.Error("failed fetching queues from db", zap.Any("error", err))
			InternalServerError(w)
			return
		}

		Ok(w, queues)

	}
}
