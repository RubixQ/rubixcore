package repositories

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SystemUser represents users of the system
type SystemUser struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Fullname  string    `json:"fullname" db:"fullname"`
	Password  string    `json:"-" db:"password"`
	IsAdmin   bool      `json:"isAdmin" db:"is_admin"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	CreatedBy int       `json:"createdBy" db:"created_by"`
	UpdatedAt time.Time `json:"updated_by" db:"updated_at"`
}

// UserRepo defines methods for interacting with db
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo returns a pointer to an initialzed
// UserRepo ready for use
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// FindByCredentials returns a system user whose login credentials
// matches what is supplied
func (r *UserRepo) FindByCredentials(username, password string) (*SystemUser, error) {
	query := "SELECT u.* FROM system_users AS u WHERE u.username = $1;"

	row := r.db.QueryRowx(query, username)
	user := new(SystemUser)

	err := row.StructScan(user)
	if err != nil {
		return nil, fmt.Errorf("failed finding user by credentials : %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("failed comparing stored password against specified password : %v", err)
	}

	if err = row.Err(); err != nil {
		return nil, fmt.Errorf(": %v", err)
	}

	return user, nil
}
