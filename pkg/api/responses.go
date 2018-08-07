package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// RenderOk is a helper for sending api response
func RenderOk(w http.ResponseWriter, data interface{}) {
	_ = renderJSON(w, data)
}

// RenderBadRequest is a helpfer for sending api response
func RenderBadRequest(w http.ResponseWriter, err error) {
	http.Error(w, fmt.Sprintf("{'error': %s}", err.Error()), http.StatusBadRequest)
	return
}

// InternalServerError is a helper for sending api response
func InternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func renderJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// response := struct {
	// 	Data interface{} `json:"data"`
	// }{
	// 	data,
	// }

	return json.NewEncoder(w).Encode(data)
}

// WriteToConn sends payload over ws connection
func WriteToConn(conn *websocket.Conn, payload WSPayload) {
	conn.WriteJSON(payload)
}
