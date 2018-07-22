package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

	a.nextTicket = a.nextTicket + 1
	customer.TicketNumber = fmt.Sprintf("%s%03d", []string{"A", "B", "C"}[rand.Intn(3)], a.nextTicket)
	customer, err = repo.Create(customer)
	if err != nil {
		a.logger.Error("failed inserting customer", zap.Error(err))
		return
	}

	go func() {
		msg := fmt.Sprintf("Your ticket number is %s. Kindly wait patiently until your turn is announced!", customer.TicketNumber)
		sendSMS(msg, customer.MSISDN)
	}()

	Ok(w, customer)
}
