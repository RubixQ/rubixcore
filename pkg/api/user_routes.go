package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/rubixq/rubixcore/pkg/db"
)

func (a *App) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := new(db.User)

		err := json.NewDecoder(r.Body).Decode(u)
		if err != nil {
			a.logger.Error("failed decoding request body into user", zap.Error(err))
			RenderBadRequest(w, err)
			return
		}

		session := a.session.Copy()
		defer session.Close()

		repo := db.NewUserRepo(session)
		u, err = repo.Create(u)
		if err != nil {
			a.logger.Error("failed inserting user into db", zap.Error(err))
			InternalServerError(w)
			return
		}

		RenderOk(w, u)
	}
}

func (a *App) listUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := a.session.Copy()
		defer session.Close()

		repo := db.NewUserRepo(session)
		users, err := repo.FindAll()
		if err != nil {
			a.logger.Error("failed fetching all users from db", zap.Error(err))
			InternalServerError(w)
			return
		}

		RenderOk(w, users)
	}
}
