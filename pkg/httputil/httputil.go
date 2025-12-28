package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/ciameksw/mood-api/pkg/logger"
)

// Helper function to handle errors
func HandleError(logger logger.Logger, w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		logger.Error.Printf("%s: %v", message, err)
	} else {
		logger.Error.Println(message)
	}
	http.Error(w, message, statusCode)
}

// Helper function to write JSON responses
func WriteJSON(logger logger.Logger, w http.ResponseWriter, data interface{}, statusCode int) {
	j, err := json.Marshal(data)
	if err != nil {
		HandleError(logger, w, "Failed to encode response to JSON", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}
