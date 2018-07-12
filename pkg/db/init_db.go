package db

import (
	"gopkg.in/mgo.v2"
)

const (
	databaseName         string = "rubixcore"
	queuesCollectionName string = "queues"
)

// InitDB sets db contraints and indexes
func InitDB(s *mgo.Session) error {
	session := s.Copy()
	defer session.Close()

	c := session.DB(databaseName).C(queuesCollectionName)
	index := mgo.Index{
		Key:    []string{"name"},
		Unique: true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		return err
	}

	return nil
}
