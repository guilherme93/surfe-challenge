package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"surfe-actions/internal/actions"
	actionshandler "surfe-actions/internal/api/actions"
	usershandler "surfe-actions/internal/api/users"
	"surfe-actions/internal/users"
)

func NewRouter(userService users.Service, actionsService actions.Service) http.Handler {
	router := chi.NewRouter()

	router.Route(
		"/api/v1", func(r chi.Router) {
			userHandler := usershandler.NewUserHandler(userService)
			actionsHanler := actionshandler.NewActionsHandler(actionsService)

			r.Route(
				"/users", func(r chi.Router) {
					r.Get("/{userId}", userHandler.GetUserByID())
				},
			)

			r.Route(
				"/actions", func(r chi.Router) {
					r.Get("/count", actionsHanler.Count())
					r.Get("/{actionType}/next-actions", actionsHanler.NextActions())
					r.Get("/referrals", actionsHanler.GetReferrals())
				},
			)
		},
	)

	return router
}
