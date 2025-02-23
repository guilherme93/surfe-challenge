package actions

import (
	"sort"

	"surfe-actions/internal/actions/domain"
	"surfe-actions/internal/actions/repository"
	"surfe-actions/internal/utils"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=actions
type Service interface {
	PredictNextActions(currentAction string) []domain.Prediction
	CountActionsByUser(userID int) int
	GetReferrals() map[int]int
}

type service struct {
	storage          repository.Repository
	nextActionMatrix map[string]map[string]int
	totalActions     map[string]int
}

func NewService(s repository.Repository) Service {
	return &service{
		storage:          s,
		nextActionMatrix: make(map[string]map[string]int),
		totalActions:     make(map[string]int),
	}
}

// PredictNextActions the first time this method is called it will be slower since we have to calculate the predictions
// a way to improve this would be to have a cron job to calculate this from time to time, so that we dont have the hit
// on the first request
func (s *service) PredictNextActions(currentAction string) []domain.Prediction {
	const roundingFloats = 2

	if len(s.totalActions) == 0 {
		s.calculatePredictions()
	}

	predictionToOrder := make([]domain.Prediction, 0)

	if _, exists := s.nextActionMatrix[currentAction]; !exists {
		return predictionToOrder
	}

	totalTransitions := s.totalActions[currentAction]
	if totalTransitions == 0 {
		return predictionToOrder
	}

	for nextAction, count := range s.nextActionMatrix[currentAction] {
		probability := float64(count) / float64(totalTransitions)

		predictionToOrder = append(predictionToOrder, domain.Prediction{
			Action:      nextAction,
			Probability: utils.RoundTo(probability, roundingFloats),
		})
	}

	sort.Slice(predictionToOrder, func(i, j int) bool {
		return predictionToOrder[i].Probability > predictionToOrder[j].Probability
	})

	return predictionToOrder
}

func (s *service) calculatePredictions() {
	actions := s.storage.GetAll()

	sort.Slice(actions, func(i, j int) bool {
		return actions[i].CreatedAt.Before(actions[j].CreatedAt)
	})

	for i := range len(actions) - 1 {
		currentAction := actions[i].Type
		nextAction := actions[i+1].Type

		if _, exists := s.nextActionMatrix[currentAction]; !exists {
			s.nextActionMatrix[currentAction] = make(map[string]int)
		}

		s.nextActionMatrix[currentAction][nextAction]++
		s.totalActions[currentAction]++
	}
}

func (s *service) CountActionsByUser(userID int) int {
	return s.storage.CountActionsByUser(userID)
}

func (s *service) GetReferrals() map[int]int {
	referrals := s.storage.GetReferrals()

	allUsers := make(map[int]struct{})
	graph := make(map[int][]int)
	memo := make(map[int]int)

	for _, ref := range referrals {
		graph[ref.UserID] = append(graph[ref.UserID], ref.TargetUser)
		allUsers[ref.UserID] = struct{}{}
		allUsers[ref.TargetUser] = struct{}{}
	}

	result := make(map[int]int)
	for user := range allUsers {
		result[user] = dfsWithMemo(user, graph, memo)
	}

	return result
}

// dfs with memoization the worst case scenario is O(N+E)
// each node is visited once due to memoization, and each edge is processed once
func dfsWithMemo(user int, referrals map[int][]int, memo map[int]int) int {
	if val, ok := memo[user]; ok {
		return val
	}

	count := 0

	for _, child := range referrals[user] {
		if child != user { // referenced himself? I am counting it as 0
			count += 1 + dfsWithMemo(child, referrals, memo)
		}
	}

	memo[user] = count

	return count
}

// this was my first try, a naive solution with O(n²) worst-case
// this case is where there is a continuous long chain of references
//
//nolint:unused
func dfsForEach(user int, referrals map[int][]int) int {
	count := 0

	for _, child := range referrals[user] {
		count += 1 + dfsForEach(child, referrals)
	}

	return count
}
