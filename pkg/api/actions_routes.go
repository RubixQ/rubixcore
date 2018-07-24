package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rubixq/rubixcore/pkg/db"

	"go.uber.org/zap"
)

func (a *App) callNextCustomer(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		QueueID   string `json:"queueID"`
		CounterID string `json:"counterID"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		a.logger.Error("failed decoding request payload", zap.Error(err))
		return
	}

	data, err := a.redis.RPop(payload.QueueID).Result()
	if err != nil {
		a.logger.Error("failed getting next ticket number", zap.Error(err))
		return
	}

	customer := db.Customer{}
	err = json.NewDecoder(strings.NewReader(data)).Decode(&customer)
	if err != nil {
		a.logger.Error("failed deserializing customer from queue", zap.Error(err))
		return
	}

	msg := fmt.Sprintf("Ticket number %s, please proceed to counter number %s", customer.TicketNumber, payload.CounterID)
	conn, ok := a.counters[payload.CounterID]
	if ok {
		payload := WSPayload{PayloadType: "update", Data: customer.TicketNumber}
		a.logger.Info("sending ws update", zap.Any("response", payload))
		WriteToConn(conn, payload)
	}

	go func() {
		sendSMS(msg, customer.MSISDN)
		a.logger.Info("next command", zap.String("command", msg))
	}()

	Ok(w, struct {
		QueueID    string `json:"queueID"`
		CounterID  string `json:"counterID"`
		NextTicket string `json:"nextTicket"`
	}{
		QueueID:    payload.QueueID,
		CounterID:  payload.CounterID,
		NextTicket: customer.TicketNumber,
	})

}
