package repo

import (
	"strings"
	"time"

	"github.com/rubixq/rubixcore/pkg/db"
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
		database:   RubixDatabase,
		collection: QueuesCollection,
	}

	qr.session = s

	return &qr
}

// Create persists a new Queue into the db
func (r *QueueRepo) Create(q *db.Queue) (*db.Queue, error) {
	q.ID = bson.NewObjectId()
	q.Title = strings.ToUpper(strings.Replace(q.Name, " ", "", -1))
	q.Active = true
	q.CreatedAt = time.Now()

	err := r.session.DB(r.database).C(r.collection).Insert(q)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// FindAll returns a list of all queues from the db
func (r *QueueRepo) FindAll() ([]db.Queue, error) {
	var queues []db.Queue

	err := r.session.DB(r.database).C(r.collection).Find(nil).All(&queues)
	if err != nil {
		return nil, err
	}

	return queues, nil
}
