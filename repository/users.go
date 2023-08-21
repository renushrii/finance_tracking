package repository

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/jmoiron/sqlx"
)

var (
	ErrMissingData = errors.New("missing data, one or more fields are empty")
)

type Users interface {
	FindByEmail(email string) (User, error)
	CreateUser(u User) error
}

type User struct {
	ID        int       `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Password  string    `db:"password" json:"password"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Validate checks if the user struct is valid
func (u *User) Validate() error {
	switch "" {
	case u.FirstName, u.LastName, u.Email, u.Password:
		return ErrMissingData
	}

	if !govalidator.IsEmail(u.Email) {
		return errors.New("invalid email")
	}

	return nil
}

type users struct {
	db *sqlx.DB
}

func NewUsers(db *sqlx.DB) Users {
	return &users{
		db: db,
	}
}

func (u *users) FindByEmail(email string) (User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = ? LIMIT 1`
	err := u.db.Get(&user, query, email)
	return user, err
}

func (u *users) CreateUser(user User) error {
	query := `INSERT INTO users (first_name,last_name, password, email) VALUES (:first_name, :last_name, :password, :email)`
	_, err := u.db.NamedExec(query, user)
	return err
}
