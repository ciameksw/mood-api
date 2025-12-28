package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ciameksw/mood-api/advice/internal/advice/repository"
	"github.com/ciameksw/mood-api/pkg/httputil"
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
		httputil.HandleError(*s.Logger, w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Var(input, "required,dive")
	if err != nil {
		httputil.HandleError(*s.Logger, w, err.Error(), err, http.StatusBadRequest)
		return
	}

	moodSummary := convertToMoodSummary(input)

	adviceTypeID, err := s.DBOperations.GetAdviceTypeIDByMoodSummary(moodSummary)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to get advice type ID", err, http.StatusInternalServerError)
		return
	}

	adviceID, title, content, err := s.DBOperations.SelectRandomAdviceByAdviceTypeID(adviceTypeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.HandleError(*s.Logger, w, "No advice found", err, http.StatusNoContent)
			return
		}
		httputil.HandleError(*s.Logger, w, "Failed to select advice", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"adviceId": adviceID,
		"title":    title,
		"content":  content,
	}
	httputil.WriteJSON(*s.Logger, w, response, http.StatusOK)
}

// Helper function to convert input to MoodSummaryEntry slice
func convertToMoodSummary(input []selectAdviceInputEntry) []repository.MoodSummaryEntry {
	summary := make([]repository.MoodSummaryEntry, len(input))
	for i, entry := range input {
		summary[i] = repository.MoodSummaryEntry{
			MoodTypeID: entry.MoodTypeID,
			Percentage: entry.Percentage,
		}
	}
	return summary
}
