package repository

import "strings"

//go:generate mockgen -source=repository.go -destination=repository_mock.go -package=repository
type Repository interface {
	CountActionsByUser(userID int) int
	GetAll() []Actions
	GetReferrals() []Actions
}

type repository struct {
	actions []Actions
}

func NewRepository() (Repository, error) {
	var err error

	actions, err := readAndParseActions()
	if err != nil {
		return nil, err
	}

	return repository{actions: actions}, nil
}

func (r repository) CountActionsByUser(userID int) int {
	var count int

	for _, action := range r.actions {
		if action.UserID == userID {
			count++
		}
	}

	return count
}

func (r repository) GetAll() []Actions {
	return r.actions
}

func (r repository) GetReferrals() []Actions {
	actions := make([]Actions, 0)

	for _, action := range r.actions {
		if strings.EqualFold(action.Type, "REFER_USER") {
			actions = append(actions, action)
		}
	}

	return actions
}
