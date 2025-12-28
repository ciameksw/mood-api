package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetTodayQuote(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting today's quote")

	resp, err := s.ExternalQuotesService.GetTodayQuote()
	if err != nil {
		s.handleError(w, "Failed to get today's quote", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, resp, http.StatusOK)
}

// Helper function to handle errors
func (s *Server) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		s.Logger.Error.Printf("%s: %v", message, err)
	} else {
		s.Logger.Error.Println(message)
	}
	http.Error(w, message, statusCode)
}

// Helper function to write JSON responses
func (s *Server) writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	j, err := json.Marshal(data)
	if err != nil {
		s.handleError(w, "Failed to encode response to JSON", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(j)
}
