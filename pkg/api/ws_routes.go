package api

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader websocket.Upgrader

func init() {
	upgrader = websocket.Upgrader{}
}

func handleStatusCheck(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
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
			logger.Error("failed upgrading request to ws connection", zap.Error(err))
		}
	}
}

func handleStatusTest(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("ws test page accessed", zap.Time("at", time.Now()))
		http.ServeFile(w, r, "ws-test.html")
	}
}
