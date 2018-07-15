package api

import (
	"html/template"
	"net/http"
	"time"

	"github.com/rubixq/rubixcore/pkg/db"
	"go.uber.org/zap"
)

func (a *App) handleKioskRegistration(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Error("failed upgrading new kiosk connection to websocket", zap.Error(err))
		return
	}

	a.addKiosk(id, conn)

	session := a.session.Copy()
	defer session.Close()
	repo := db.NewQueueRepo(session)

	queues, err := repo.FindAll()
	if err != nil {
		conn.WriteJSON(
			wsPayload{
				ptype: "error",
				data:  err,
			})
	}

	conn.WriteJSON(
		wsPayload{
			ptype: "queues",
			data:  queues,
		})
}

func (a *App) handleKioskSetup(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("ws test page accessed", zap.Time("at", time.Now()))

	files := []string{
		"./ui/html/base.html",
		"./ui/html/kiosk.page.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.logger.Error("failed parsing templates", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		a.logger.Error("failed executing template", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
