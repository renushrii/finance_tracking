package server

import (
	"testing"

	"github.com/renushrii/finance-tracking/repository"
)

func TestHashFn(t *testing.T) {
	password := "123456"

	hashedPassword, err := hash(password)
	if err != nil {
		t.Error(err)
	}

	if !compareHash(hashedPassword, password) {
		t.Errorf("Hashed password does not match: hashed %s, original %s", hashedPassword, password)
	}
}

func TestLoginAPI(t *testing.T) {
	s := &Server{
		users: &MockUsers{users: []repository.User{}},
	}

	s.users.CreateUser(repository.User{
		FirstName: "Renushri",
		LastName:  "Rawat",
		Email:     "renu@gmail.com",
		Password:  "1234",
	})

	if _, err := s.login("renu@gmail.com", "1234"); err != nil {
		t.Error(err)
	}

	if _, err := s.login("renu@gmail.com", "qwerty"); err == nil {
		t.Error("Should have failed")
	}

	if _, err := s.login("some@gmail.com", "234567"); err == nil {
		t.Error("Should have failed")
	}
}

type MockUsers struct {
	users []repository.User
}

func (m *MockUsers) FindByEmail(email string) (repository.User, error) {
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return repository.User{}, nil
}

func (m *MockUsers) CreateUser(user repository.User) error {
	hashed, err := hash(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashed
	user.ID = len(m.users) + 1

	m.users = append(m.users, user)
	return nil
}
