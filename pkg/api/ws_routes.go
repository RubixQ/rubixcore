package api

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (a *App) handleStatusCheck(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err == nil {
		go func(c *websocket.Conn) {
			ticker := time.NewTicker(10 * time.Second)
			for now := range ticker.C {
				c.WriteJSON(
					struct {
						Message string    `json:"msg"`
						At      time.Time `json:"at"`
					}{
						"Websocket Status Message",
						now,
					})
			}
		}(conn)
	} else {
		a.logger.Error("failed upgrading request to ws connection", zap.Error(err))
	}
}

func (a *App) handleStatusTest(w http.ResponseWriter, r *http.Request) {
	a.logger.Info("ws test page accessed", zap.Time("at", time.Now()))

	files := []string{
		"./ui/html/base.html",
		"./ui/html/ws.page.html",
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
