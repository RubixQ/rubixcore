package db

import (
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CustomerRepo defines methods for managing queues in the db
type CustomerRepo struct {
	database   string
	collection string
	session    *mgo.Session
}

// NewCustomerRepo constructs and returns a pointer to a new CustomerRepo
func NewCustomerRepo(s *mgo.Session) *CustomerRepo {
	qr := CustomerRepo{
		database:   databaseName,
		collection: customersCollectionName,
	}

	qr.session = s

	return &qr
}

// Create persists a new Customer into the db
func (r *CustomerRepo) Create(c *Customer) (*Customer, error) {
	c.ID = bson.NewObjectId()
	c.CreatedAt = time.Now()

	err := r.session.DB(r.database).C(r.collection).Insert(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
