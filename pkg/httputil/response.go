package httputil

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Status  int         `json:"-"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// JSON sends a JSON response with the appropriate status code
func JSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

// SendSuccess sends a successful response
func SendSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	resp := Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}
	JSON(w, statusCode, resp)
}

// SendError sends an error response
func SendError(w http.ResponseWriter, statusCode int, err string) {
	resp := Response{
		Status: statusCode,
		Error:  err,
	}
	JSON(w, statusCode, resp)
}
