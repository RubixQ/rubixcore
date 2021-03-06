package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/rubixq/rubixcore/pkg/db"
	"go.uber.org/zap"
)

func (a *App) createCustomer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := new(db.Customer)

		err := json.NewDecoder(r.Body).Decode(&customer)
		if err != nil {
			a.logger.Error("failed decoding request payload", zap.Error(err))
			failedCustomerRegistrations.Inc()

			RenderBadRequest(w, err)
		}

		session := a.session.Copy()
		defer session.Close()

		repo := db.NewCustomerRepo(session)

		a.nextTicket = a.nextTicket + 1
		customer.TicketNumber = fmt.Sprintf("%s%03d", []string{"A", "B", "C", "X", "Y", "Z"}[rand.Intn(6)], a.nextTicket)
		customer, err = repo.Create(customer)
		if err != nil {
			a.logger.Error("failed inserting customer", zap.Error(err))
			failedCustomerRegistrations.Inc()

			RenderBadRequest(w, err)
		}

		data, err := json.Marshal(customer)
		if err != nil {
			a.logger.Error("failed serializing customer data", zap.Error(err))
			return
		}
		a.redis.LPush(customer.QueueID, string(data))

		go func() {
			msg := fmt.Sprintf("Your ticket number is %s. Kindly wait patiently until your turn is announced!", customer.TicketNumber)
			sendSMS(msg, customer.MSISDN)
		}()

		RenderOk(w, customer)
	}
}

func (a *App) listCustomers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := a.session.Copy()
		defer session.Close()

		repo := db.NewCustomerRepo(session)

		customers, err := repo.FindAll()

		if err != nil {
			a.logger.Error("failed fetching customers from db", zap.Any("error", err))
			InternalServerError(w)
			return
		}

		RenderOk(w, customers)

	}
}
