package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ciameksw/mood-api/advice/internal/advice/postgres"
)

type selectAdviceInputEntry struct {
	MoodTypeID int     `json:"moodTypeId" validate:"required"`
	Count      int     `json:"count" validate:"required,min=1"`
	Percentage float64 `json:"percentage" validate:"required"`
}

func (s *Server) handleSelectAdvice(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Selecting advice")
	var input []selectAdviceInputEntry

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Var(input, "required,dive")
	if err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	moodSummary := convertToMoodSummary(input)

	adviceTypeID, err := s.Postgres.GetAdviceTypeIDByMoodSummary(moodSummary)
	if err != nil {
		s.handleError(w, "Failed to get advice type ID", err, http.StatusInternalServerError)
		return
	}

	adviceID, title, content, err := s.Postgres.SelectRandomAdviceByAdviceTypeID(adviceTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.handleError(w, "No advice found", err, http.StatusNoContent)
			return
		}
		s.handleError(w, "Failed to select advice", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"adviceId": adviceID,
		"title":    title,
		"content":  content,
	}
	s.writeJSON(w, response, http.StatusOK)
}

// Helper function to convert input to MoodSummaryEntry slice
func convertToMoodSummary(input []selectAdviceInputEntry) []postgres.MoodSummaryEntry {
	summary := make([]postgres.MoodSummaryEntry, len(input))
	for i, entry := range input {
		summary[i] = postgres.MoodSummaryEntry{
			MoodTypeID: entry.MoodTypeID,
			Percentage: entry.Percentage,
		}
	}
	return summary
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
