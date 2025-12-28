package server

import (
	"errors"
	"io"
	"net/http"
	"time"
)

// Helper function to handle errors
func (s *Server) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		s.Logger.Error.Printf("%s: %v", message, err)
	} else {
		s.Logger.Error.Println(message)
	}
	http.Error(w, message, statusCode)
}

// Helper function to forward the response
func (s *Server) forwardResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// Helper function to parse timeframe query parameters for mood and advice retrieval
func (s *Server) parseQueryParams(r *http.Request) (string, string, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		return "", "", errors.New("from and to parameters are required")
	}

	// Validate date format YYYY-MM-DD
	const dateFormat = "2006-01-02"
	if _, err := time.Parse(dateFormat, from); err != nil {
		return "", "", errors.New("from date must be in YYYY-MM-DD format")
	}
	if _, err := time.Parse(dateFormat, to); err != nil {
		return "", "", errors.New("to date must be in YYYY-MM-DD format")
	}

	return from, to, nil
}
