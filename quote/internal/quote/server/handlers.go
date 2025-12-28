package server

import (
	"context"
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetTodayQuote(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting today's quote")

	ctx := context.Background()

	cachedQuote, err := s.RedisCache.GetTodayQuote(ctx)
	if err == nil && cachedQuote != nil {
		s.Logger.Info.Println("Quote found in cache")
		s.writeJSON(w, cachedQuote, http.StatusOK)
		return
	}

	resp, err := s.ExternalQuotesService.GetTodayQuote()
	if err != nil {
		s.handleError(w, "Failed to get today's quote", err, http.StatusInternalServerError)
		return
	}

	if err := s.RedisCache.SetTodayQuote(ctx, resp); err != nil {
		s.Logger.Error.Printf("Failed to cache quote: %v", err)
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
