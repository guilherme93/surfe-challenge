package actions

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"surfe-actions/internal/actions"
	"surfe-actions/internal/api"
	"surfe-actions/internal/api/actions/dto"
)

type Handler interface {
	Count() http.HandlerFunc
	NextActions() http.HandlerFunc
	GetReferrals() http.HandlerFunc
}

type actionsHandler struct {
	service actions.Service
}

func NewActionsHandler(service actions.Service) Handler {
	return actionsHandler{service: service}
}

func (a actionsHandler) Count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIDQuery := r.URL.Query().Get("userId")

		userID, err := strconv.Atoi(userIDQuery)
		if err != nil {
			api.RenderError(w, http.StatusBadRequest, errors.New("invalid user ID"))

			return
		}

		count := a.service.CountActionsByUser(userID)

		resp := dto.CountResponse{Count: count}
		api.RenderJSON(w, http.StatusOK, resp)
	}
}

func (a actionsHandler) NextActions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		actionType := chi.URLParam(r, "actionType")

		predictions := a.service.PredictNextActions(actionType)
		if len(predictions) == 0 {
			api.RenderError(w, http.StatusNotFound, errors.New(actionType+" not found"))

			return
		}

		resp := make([]dto.Prediction, 0, len(predictions))
		for _, prediction := range predictions {
			resp = append(resp, dto.Prediction{
				Action:      prediction.Action,
				Probability: prediction.Probability,
			})
		}

		api.RenderJSON(w, http.StatusOK, resp)
	}
}

func (a actionsHandler) GetReferrals() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		result := a.service.GetReferrals()

		api.RenderJSON(w, http.StatusOK, result)
	}
}
