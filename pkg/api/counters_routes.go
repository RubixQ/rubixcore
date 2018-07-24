package api

import (
	"html/template"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (a *App) handleCounterSetup(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("counter test page accessed", zap.Time("at", time.Now()))

	files := []string{
		"./ui/html/base.html",
		"./ui/html/counter.page.html",
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

func (a *App) handleCounterWebsocket(w http.ResponseWriter, r *http.Request) {
	payload := CounterRegPayload{}
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Info("failed upgrading to ws connection", zap.Error(err))
		return
	}

	err = conn.ReadJSON(&payload)
	if err != nil {
		a.logger.Error("failed reading counter registration payload", zap.Error(err))
	}

	a.addCounter(payload.CounterID, conn)

	conn.WriteJSON(
		WSPayload{
			PayloadType: "welcome",
			Data:        "connected successfully",
		},
	)
}
