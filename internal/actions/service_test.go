package actions

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"surfe-actions/internal/actions/domain"
	actionsrepository "surfe-actions/internal/actions/repository"
)

func TestGetReferrals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		referrals []actionsrepository.Actions
		expected  map[int]int
	}{
		{
			name:      "no referrals",
			referrals: []actionsrepository.Actions{},
			expected:  map[int]int{},
		},
		{
			name: "single referral",
			referrals: []actionsrepository.Actions{
				{UserID: 1, TargetUser: 2},
			},
			expected: map[int]int{
				1: 1,
				2: 0,
			},
		},
		{
			name: "multiple referrals",
			referrals: []actionsrepository.Actions{
				{UserID: 1, TargetUser: 2},
				{UserID: 2, TargetUser: 3},
				{UserID: 3, TargetUser: 4},
			},
			expected: map[int]int{
				1: 3,
				2: 2,
				3: 1,
				4: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := actionsrepository.NewMockRepository(ctrl)

			mockStorage.EXPECT().GetReferrals().Return(tt.referrals)

			s := &service{storage: mockStorage}

			result := s.GetReferrals()

			assert.True(t, reflect.DeepEqual(result, tt.expected))
		})
	}
}

func TestPredictNextActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := actionsrepository.NewMockRepository(ctrl)

	tests := []struct {
		name          string
		currentAction string
		mockActions   []actionsrepository.Actions
		expected      []domain.Prediction
	}{
		{
			name:          "no actions",
			currentAction: "A",
			mockActions:   []actionsrepository.Actions{},
			expected:      []domain.Prediction{},
		},
		{
			name:          "single action",
			currentAction: "A",
			mockActions:   []actionsrepository.Actions{{Type: "A", CreatedAt: time.Now()}},
			expected:      []domain.Prediction{},
		},
		{
			name:          "basic transition",
			currentAction: "A",
			mockActions: []actionsrepository.Actions{
				{Type: "A", CreatedAt: time.Now()},
				{Type: "B", CreatedAt: time.Now().Add(time.Second)},
			},
			expected: []domain.Prediction{{Action: "B", Probability: 1.0}},
		},
		{
			name:          "multiple transitions",
			currentAction: "A",
			mockActions: []actionsrepository.Actions{
				{Type: "A", CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Type: "B", CreatedAt: time.Date(2023, 1, 1, 0, 0, 1, 0, time.UTC)},
				{Type: "A", CreatedAt: time.Date(2023, 1, 1, 0, 0, 2, 0, time.UTC)},
				{Type: "C", CreatedAt: time.Date(2023, 1, 1, 0, 0, 3, 0, time.UTC)},
				{Type: "A", CreatedAt: time.Date(2023, 1, 1, 0, 0, 4, 0, time.UTC)},
				{Type: "C", CreatedAt: time.Date(2023, 1, 1, 0, 0, 4, 0, time.UTC)},
				{Type: "A", CreatedAt: time.Date(2023, 1, 1, 0, 0, 4, 0, time.UTC)},
			},
			expected: []domain.Prediction{
				{Action: "C", Probability: 0.67},
				{Action: "B", Probability: 0.33},
			},
		},
		{
			name:          "action not in matrix",
			currentAction: "D",
			mockActions: []actionsrepository.Actions{
				{Type: "A", CreatedAt: time.Now()},
				{Type: "B", CreatedAt: time.Now().Add(time.Second)},
			},
			expected: []domain.Prediction{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				storage:          mockStorage,
				nextActionMatrix: make(map[string]map[string]int),
				totalActions:     make(map[string]int),
			}

			mockStorage.EXPECT().GetAll().Return(tt.mockActions)

			result := svc.PredictNextActions(tt.currentAction)

			assert.Equal(t, tt.expected, result)
		})
	}
}
