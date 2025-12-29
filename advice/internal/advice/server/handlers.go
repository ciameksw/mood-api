package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

	id, err := s.DBOperations.SaveUserAdvicePeriod(req.UserID, req.AdviceID, req.From, req.To)
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

func (s *Server) handleGetByID(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting advice by ID")

	idStr := r.PathValue("id")
	if idStr == "" {
		httputil.HandleError(*s.Logger, w, "Missing id parameter", nil, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid id parameter", err, http.StatusBadRequest)
		return
	}

	title, content, err := s.DBOperations.GetAdviceByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.HandleError(*s.Logger, w, "Advice not found", err, http.StatusNotFound)
			return
		}
		httputil.HandleError(*s.Logger, w, "Failed to get advice by ID", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"adviceId": id,
		"title":    title,
		"content":  content,
	}
	httputil.WriteJSON(*s.Logger, w, response, http.StatusOK)
}

func (s *Server) handleGetAdviceByPeriod(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting advice for period")

	input, err := s.parseQueryParams(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	adviceID, title, content, err := s.DBOperations.GetAdviceByPeriod(input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httputil.HandleError(*s.Logger, w, "Advice not found for given period", err, http.StatusNotFound)
			return
		}
		httputil.HandleError(*s.Logger, w, "Failed to get advice for period", err, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"adviceId": adviceID,
		"title":    title,
		"content":  content,
	}
	httputil.WriteJSON(*s.Logger, w, response, http.StatusOK)
}

type GetInput struct {
	UserID    int
	StartDate string
	EndDate   string
}

func (s *Server) parseQueryParams(r *http.Request) (*GetInput, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		return nil, err
	}

	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		return nil, errors.New("from and to parameters are required")
	}

	// Validate date format YYYY-MM-DD
	const dateFormat = "2006-01-02"
	if _, err := time.Parse(dateFormat, from); err != nil {
		return nil, errors.New("from date must be in YYYY-MM-DD format")
	}
	if _, err := time.Parse(dateFormat, to); err != nil {
		return nil, errors.New("to date must be in YYYY-MM-DD format")
	}

	return &GetInput{
		UserID:    userID,
		StartDate: from,
		EndDate:   to,
	}, nil
}
