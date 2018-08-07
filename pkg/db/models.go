package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// User defines model for system users
type User struct {
	ID        bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	Username  string        `bson:"username,omitempty" json:"username,omitempty"`
	Password  string        `bson:"password,omitempty" json:"-"`
	IsAdmin   bool          `bson:"isAdmin,omitempty" json:"isAdmin,omitempty"`
	IsActive  bool          `bson:"isActive,omitempty" json:"isActive,omitempty"`
	CreatedAt time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time     `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// Queue defines model for queues
type Queue struct {
	ID          bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Title       string        `bson:"title,omitempty" json:"title,omitempty"`
	IsActive    bool          `bson:"isActive,omitempty" json:"isActive,omitempty"`
	CreatedAt   time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time     `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// Customer defines model for customers
type Customer struct {
	ID           bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	MSISDN       string        `bson:"msisdn,omitempty" json:"msisdn,omitempty"`
	QueueID      string        `bson:"queueId,omitempty" json:"queueId,omitempty"`
	TicketNumber string        `bson:"ticketNumber,omitempty" json:"ticketNumber,omitempty"`
	IsServed     bool          `bson:"isServed,omitempty" json:"isServed,omitempty"`
	ServedAt     time.Time     `bson:"servedAt,omitempty" json:"servedAt,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt    time.Time     `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
