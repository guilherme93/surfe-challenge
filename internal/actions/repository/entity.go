package repository

import (
	"embed"
	"encoding/json"
	"time"
)

type Actions struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	UserID     int       `json:"userID"`
	TargetUser int       `json:"targetUser"`
	CreatedAt  time.Time `json:"createdAt"`
}

//go:embed  "actions.json"
var embeddedFiles embed.FS

func readAndParseActions() ([]Actions, error) {
	result, err := embeddedFiles.ReadFile("actions.json")
	if err != nil {
		return nil, err
	}

	actions := make([]Actions, 0)

	err = json.Unmarshal(result, &actions)
	if err != nil {
		return nil, err
	}

	return actions, nil
}
