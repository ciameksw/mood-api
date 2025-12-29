package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

type saveAdviceRequest struct {
	UserID   int    `json:"userId" validate:"required"`
	AdviceID int    `json:"adviceId" validate:"required"`
	From     string `json:"from" validate:"required,datetime=2006-01-02"`
	To       string `json:"to" validate:"required,datetime=2006-01-02"`
}

func (s *Server) handleSaveAdvice(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Saving advice period")

	var req saveAdviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	if err := s.Validator.Struct(req); err != nil {
		httputil.HandleError(*s.Logger, w, err.Error(), err, http.StatusBadRequest)
		return
	}

	parsedFrom, err := time.Parse("2006-01-02", req.From)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid 'from' date format, expected YYYY-MM-DD", err, http.StatusBadRequest)
		return
	}

	parsedTo, err := time.Parse("2006-01-02", req.To)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid 'to' date format, expected YYYY-MM-DD", err, http.StatusBadRequest)
		return
	}

	if parsedFrom.After(parsedTo) {
		httputil.HandleError(*s.Logger, w, "'from' must be before or equal to 'to'", errors.New("invalid date range"), http.StatusBadRequest)
		return
	}

	id, err := s.DBOperations.SaveUserAdvicePeriod(req.UserID, req.AdviceID, parsedFrom, parsedTo)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to save advice period", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":         id,
		"userId":     req.UserID,
		"adviceId":   req.AdviceID,
		"periodFrom": req.From,
		"periodTo":   req.To,
	}
	httputil.WriteJSON(*s.Logger, w, response, http.StatusCreated)
}
