package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (a *App) handleVoiceWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.logger.Error("failed upgrading to ws connection", zap.Error(err))
		return
	}

	a.voiceSvr = conn

	conn.WriteJSON(
		WSPayload{
			PayloadType: "welcome",
			Data:        "connected successfully",
		},
	)
}
