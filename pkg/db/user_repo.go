package db

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserRepo define methods for managin users in the db
type UserRepo struct {
	database   string
	collection string
	session    *mgo.Session
}

// NewUserRepo returns a pointer to a new UserRepo
func NewUserRepo(s *mgo.Session) *UserRepo {
	ur := UserRepo{
		database:   databaseName,
		collection: usersCollectionName,
		session:    s,
	}

	return &ur
}

// Create persists a new User into the db
func (r *UserRepo) Create(u *User) (*User, error) {
	u.ID = bson.NewObjectId()
	u.IsActive = true
	u.CreatedAt = time.Now()

	err := r.session.DB(r.database).C(r.collection).Insert(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// FindAll returns a list of all users from the db
func (r *UserRepo) FindAll() ([]User, error) {
	var users []User

	err := r.session.DB(r.database).C(r.collection).Find(nil).All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// FindByCredentials returns a user whose username and password matches
func (r *UserRepo) FindByCredentials(username, password string) (*User, error) {
	var u User

	err := r.session.DB(r.database).C(r.collection).Find(bson.M{"username": username, "password": password}).One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
