package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

// Ok is a helper for sending api response
func Ok(w http.ResponseWriter, data interface{}) {
	_ = renderJSON(w, data)
}

// BadRequest is a helpfer for sending api response
func BadRequest(w http.ResponseWriter) {

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
