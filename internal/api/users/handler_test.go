package users_test

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	usershandler "surfe-actions/internal/api/users"
	"surfe-actions/internal/users"
	"surfe-actions/internal/users/domain"
	"testing"
	"time"
)

func TestGetUserByIDHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		userID         string
		mock           func(service *users.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "valid user ID",
			userID: "123",
			mock: func(service *users.MockService) {
				service.EXPECT().
					GetUserByID(123).
					Return(
						&domain.User{
							ID:        123,
							Name:      "John Doe",
							CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						},
						nil,
					)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":123,"name":"John Doe","createdAt":"2023-01-01T00:00:00Z"}`,
		},
		{
			name:           "invalid user ID",
			userID:         "abc",
			mock:           func(service *users.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid user ID"}`,
		},
		{
			name:   "user not found",
			userID: "456",
			mock: func(service *users.MockService) {
				service.EXPECT().
					GetUserByID(456).
					Return(nil, users.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		},
		{
			name:   "unexpected error",
			userID: "789",
			mock: func(service *users.MockService) {
				service.EXPECT().
					GetUserByID(789).
					Return(nil, errors.New("unexpected error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"unexpected error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUsersService := users.NewMockService(ctrl)

			tt.mock(mockUsersService)

			r := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("userId", tt.userID)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := usershandler.NewUserHandler(mockUsersService)
			handler.GetUserByID()(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
