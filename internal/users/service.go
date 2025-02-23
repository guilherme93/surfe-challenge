package users

import (
	"errors"

	"surfe-actions/internal/users/domain"
	"surfe-actions/internal/users/repository"
	"surfe-actions/internal/utils"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=users
type Service interface {
	GetUserByID(userID int) (*domain.User, error)
}

type service struct {
	storage repository.Repository
}

var ErrNotFound = errors.New("user not found")

func NewService(s repository.Repository) Service {
	return service{
		storage: s,
	}
}

func (s service) GetUserByID(userID int) (*domain.User, error) {
	user, err := s.storage.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errors.Join(err, ErrNotFound)
		}

		return nil, err
	}

	return utils.ToPtr(domain.FromEntity(user)), nil
}
