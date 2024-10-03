package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Data             interface{} `json:"data,omitempty"`             // Optional
	Message          string      `json:"message"`                    // Mandatory
	StatusCode       int         `json:"statusCode"`                 // Mandatory
	IsError          bool        `json:"isError"`                    // Mandatory
	ErrorCode        int         `json:"errorCode"`                  // Optional
	ErrorMessage     string      `json:"errorMessage"`               // Optional
	DeveloperMessage string      `json:"developerMessage,omitempty"` // Optional
	UserMessage      string      `json:"userMessage,omitempty"`      // Optional
}

func WriteResponse(w http.ResponseWriter, res interface{}) error {
	// Check if res is of type Response
	r, ok := res.(Response)
	if !ok {
		// Handle the case where res is not a Response type
		http.Error(w, "invalid response type", http.StatusInternalServerError)
		return fmt.Errorf("invalid response type, expected Response, got %T", res)
	} else {
		// Now it's safe to access fields from r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.StatusCode)

		return json.NewEncoder(w).Encode(r)
	}

}
