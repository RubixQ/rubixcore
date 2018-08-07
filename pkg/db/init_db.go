package db

import (
	"gopkg.in/mgo.v2"
)

const (
	databaseName            string = "rubixcore"
	queuesCollectionName    string = "queues"
	customersCollectionName string = "customers"
	usersCollectionName     string = "users"
)

// InitDB sets db contraints and indexes
func InitDB(s *mgo.Session, adminUsername, adminPassword string) error {
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

	c = session.DB(databaseName).C(usersCollectionName)
	index = mgo.Index{
		Key:    []string{"username"},
		Unique: true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return err
	}

	count, err := c.Count()
	if err != nil {
		return err
	}

	if count == 0 {
		u := User{
			Username: adminUsername,
			Password: adminPassword,
			IsActive: true,
			IsAdmin:  true,
		}

		err = c.Insert(u)
		if err != nil {
			return err
		}
	}

	return nil
}
