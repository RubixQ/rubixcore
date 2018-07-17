package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/rubixq/rubixcore/pkg/db"
)

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
	customer := new(db.Customer)

	err := json.NewDecoder(r.Body).Decode(customer)
	if err != nil {
		a.logger.Error("failed decoding request payload", zap.Error(err))
		return
	}

	session := a.session.Copy()
	defer session.Close()

	repo := db.NewCustomerRepo(session)

	customer, err = repo.Create(customer)
	if err != nil {
		a.logger.Error("failed inserting customer", zap.Error(err))
		return
	}

	go func() {
		msg := "Your ticket number is A125"
		sendSMS(msg, customer.MSISDN)
	}()

	Ok(w, customer)
}
