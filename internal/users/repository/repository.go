package repository

import "errors"

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=repository
type Repository interface {
	GetUserByID(id int) (*User, error)
}

type repository struct {
	users []User
}

var ErrNotFound = errors.New("entity not found")

func NewRepository() (Repository, error) {
	var err error

	users, err := readAndParseUsers()
	if err != nil {
		return nil, err
	}

	return repository{users: users}, nil
}

func (r repository) GetUserByID(id int) (*User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, ErrNotFound
}

func (r repository) GetAllUsers() []User {
	return r.users
}
