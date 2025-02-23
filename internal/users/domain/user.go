package domain

import (
	"time"

	"surfe-actions/internal/users/repository"
)

type User struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

func FromEntity(user *repository.User) User {
	return User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}
}
