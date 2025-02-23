package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"surfe-actions/internal/api"
	"surfe-actions/internal/api/users/dto"
	"surfe-actions/internal/users"
)

type UserHandler interface {
	GetUserByID() http.HandlerFunc
}

type userHandler struct {
	usersService users.Service
}

func NewUserHandler(usersService users.Service) UserHandler {
	return userHandler{
		usersService: usersService,
	}
}

func (u userHandler) GetUserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDQuery := chi.URLParam(r, "userId")

		userID, err := strconv.Atoi(userIDQuery)
		if err != nil {
			api.RenderError(w, http.StatusBadRequest, errors.New("invalid user ID"))

			return
		}

		user, err := u.usersService.GetUserByID(userID)
		if err != nil {
			if errors.Is(err, users.ErrNotFound) {
				api.RenderError(w, http.StatusNotFound, errors.New("user not found"))

				return
			}

			api.RenderError(w, http.StatusInternalServerError, errors.New("unexpected error"))

			return
		}

		resp := dto.User{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		}
		api.RenderJSON(w, http.StatusOK, resp)
	}
}
