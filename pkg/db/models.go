package db

import (
	"gopkg.in/mgo.v2/bson"
)

// Queue defines model for queues
type Queue struct {
	ID          bson.ObjectId `bson:"id,omitempty" json:"id,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	Description string        `bson:"description,omitempty" json:"description,omitempty"`
	Title       string        `bson:"title,omitempty" json:"title,omitempty"`
	Active      bool          `bson:"active,omitempty" json:"active,omitempty"`
}
