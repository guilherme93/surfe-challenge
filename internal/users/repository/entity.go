package repository

import (
	"embed"
	"encoding/json"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

//go:embed  "users.json"
var embeddedFiles embed.FS

func readAndParseUsers() ([]User, error) {
	result, err := embeddedFiles.ReadFile("users.json")
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	err = json.Unmarshal(result, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
