package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type Error struct {
	Error string `json:"error"`
}

func RenderJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set(headerContentType, contentTypeJSON)

	body, err := json.Marshal(payload)
	if err != nil {
		statusCode = http.StatusInternalServerError
		body = []byte(fmt.Sprintf(`{"error":"%s"}`, err))
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(body)
}

func RenderError(w http.ResponseWriter, statusCode int, err error) {
	RenderJSON(w, statusCode, Error{Error: err.Error()})
}
