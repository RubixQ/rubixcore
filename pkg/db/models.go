package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Queue defines model for queues
type Queue struct {
	ID          bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Title       string        `bson:"title,omitempty" json:"title,omitempty"`
	Active      bool          `bson:"active,omitempty" json:"active,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time     `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// Customer defines model for customers
type Customer struct {
	ID        bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	MSISDN    string        `bson:"msisdn,omitempty" json:"msisdn,omitempty"`
	QueueID   string        `bson:"queueId,omitempty" json:"queueId,omitempty"`
	QueueName string        `bson:"queueName,omitempty" json:"queueName,omitempty"`
	CreatedAt time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time     `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
