package actions

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"surfe-actions/internal/actions"
	"surfe-actions/internal/actions/domain"
	"testing"
)

func TestCountHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		queryParams    url.Values
		mock           func(service *actions.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "valid user ID",
			queryParams: url.Values{"userId": []string{"123"}},
			mock: func(service *actions.MockService) {
				service.EXPECT().CountActionsByUser(123).Return(5)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"count":5}`,
		},
		{
			name:           "invalid user ID",
			queryParams:    url.Values{"userId": []string{"abc"}},
			mock:           func(_ *actions.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid user ID"}`,
		},
		{
			name:           "missing user ID",
			queryParams:    url.Values{},
			mock:           func(_ *actions.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid user ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := actions.NewMockService(ctrl)

			tt.mock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/count?"+tt.queryParams.Encode(), nil)
			w := httptest.NewRecorder()

			handler := &actionsHandler{service: mockService}
			handler.Count()(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestNextActionsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		actionType     string
		mock           func(service *actions.MockService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "predictions found",
			actionType: "A",
			mock: func(service *actions.MockService) {
				service.EXPECT().PredictNextActions("A").Return([]domain.Prediction{
					{Action: "B", Probability: 0.5},
					{Action: "C", Probability: 0.3},
				})
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"action":"B","probability":0.5},{"action":"C","probability":0.3}]`,
		},
		{
			name:       "no predictions found",
			actionType: "D",
			mock: func(service *actions.MockService) {
				service.EXPECT().PredictNextActions("D").Return([]domain.Prediction{})
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"D not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := actions.NewMockService(ctrl)

			tt.mock(mockService)

			r := httptest.NewRequest(http.MethodGet, "/next-actions/"+tt.actionType, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("actionType", tt.actionType)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := &actionsHandler{service: mockService}
			handler.NextActions()(w, r)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
