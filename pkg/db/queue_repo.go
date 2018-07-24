package db

import (
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// QueueRepo defines methods for managing queues in the db
type QueueRepo struct {
	database   string
	collection string
	session    *mgo.Session
}

// NewQueueRepo constructs and returns a pointer to a new QueueRepo
func NewQueueRepo(s *mgo.Session) *QueueRepo {
	qr := QueueRepo{
		database:   databaseName,
		collection: queuesCollectionName,
	}

	qr.session = s

	return &qr
}

// Create persists a new Queue into the db
func (r *QueueRepo) Create(q *Queue) (*Queue, error) {
	q.ID = bson.NewObjectId()
	q.Title = strings.ToUpper(strings.Replace(q.Name, " ", "", -1))
	q.IsActive = true
	q.CreatedAt = time.Now()

	err := r.session.DB(r.database).C(r.collection).Insert(q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// FindAll returns a list of all queues from the db
func (r *QueueRepo) FindAll() ([]Queue, error) {
	var queues []Queue

	err := r.session.DB(r.database).C(r.collection).Find(nil).All(&queues)
	if err != nil {
		return nil, err
	}

	return queues, nil
}
