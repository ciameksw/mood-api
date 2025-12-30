package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/ciameksw/mood-api/pkg/logger"
)

// response is the standard response structure for all API endpoints.
type response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// HandleError logs the error (if present) and sends a unified error response.
func HandleError(logger logger.Logger, w http.ResponseWriter, errMsg string, err error, statusCode int) {
	if err != nil {
		logger.Error.Printf("%s: %v", errMsg, err)
	} else {
		logger.Error.Println(errMsg)
	}
	resp := response{
		Error: errMsg,
	}
	writeJSON(w, resp, statusCode)
}

// WriteSuccessMessage writes a success response with a message.
func WriteSuccessMessage(logger logger.Logger, w http.ResponseWriter, message string, statusCode int) {
	resp := response{
		Message: message,
	}
	writeJSON(w, resp, statusCode)
}

// WriteData writes a unified response with arbitrary data (unwrapped, not nested under "data" field).
func WriteData(logger logger.Logger, w http.ResponseWriter, data interface{}, statusCode int) {
	writeJSON(w, data, statusCode)
}

// writeJSON helper to serialize and write unified response.
func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	j, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}
