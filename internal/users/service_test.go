package users

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"surfe-actions/internal/users/domain"
	"surfe-actions/internal/users/repository"
	"testing"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		userID        int
		mockUser      *repository.User
		mockError     error
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:   "user found",
			userID: 1,
			mockUser: &repository.User{
				ID:   1,
				Name: "John Doe",
			},
			mockError: nil,
			expectedUser: &domain.User{
				ID:   1,
				Name: "John Doe",
			},
			expectedError: nil,
		},
		{
			name:          "user not found",
			userID:        2,
			mockUser:      nil,
			mockError:     repository.ErrNotFound,
			expectedUser:  nil,
			expectedError: errors.Join(repository.ErrNotFound, ErrNotFound),
		},
		{
			name:          "storage error",
			userID:        3,
			mockUser:      nil,
			mockError:     errors.New("storage error"),
			expectedUser:  nil,
			expectedError: errors.New("storage error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := repository.NewMockRepository(ctrl)

			svc := NewService(mockStorage)

			mockStorage.EXPECT().GetUserByID(tt.userID).Return(tt.mockUser, tt.mockError)

			user, err := svc.GetUserByID(tt.userID)
			if err != nil {
				assert.Error(t, err)
				assert.EqualError(t, tt.expectedError, err.Error(), err.Error())
			} else {
				assert.NoError(t, tt.expectedError)
			}

			assert.Equal(t, tt.expectedUser, user)
		})
	}
}
