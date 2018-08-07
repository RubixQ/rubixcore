package api

import (
	"encoding/json"
	"net/http"

	"github.com/rubixq/rubixcore/pkg/db"
	"go.uber.org/zap"
)

func (a *App) createQueue() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queue := new(db.Queue)

		err := json.NewDecoder(r.Body).Decode(queue)
		if err != nil {
			a.logger.Error("failed decoding request payload", zap.Any("error", err))
			return
		}

		session := a.session.Copy()
		defer session.Close()

		repo := db.NewQueueRepo(session)

		queue, err = repo.Create(queue)
		if err != nil {
			a.logger.Error("failed inserting queue", zap.Error(err))
			return
		}

		RenderOk(w, queue)
	}

}

func (a *App) listQueues() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := a.session.Copy()
		defer session.Close()

		repo := db.NewQueueRepo(session)

		queues, err := repo.FindAll()

		if err != nil {
			a.logger.Error("failed fetching queues from db", zap.Any("error", err))
			InternalServerError(w)
			return
		}

		RenderOk(w, queues)
	}

}
