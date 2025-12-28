package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ciameksw/mood-api/mood/internal/mood/repository"
)

type addMoodInput struct {
	UserID     int    `json:"userId" validate:"required"`
	MoodTypeID int    `json:"moodTypeId" validate:"required"`
	Note       string `json:"note" validate:"max=500"`
	Date       string `json:"date" validate:"required,datetime=2006-01-02"`
}

func (s *Server) handleAddMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Adding mood entry")
	var input addMoodInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Struct(input)
	if err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	entry, err := s.DBOperations.GetMoodEntryByDateAndUser(input.UserID, input.Date)
	if err != nil && err.Error() != "user not found" {
		s.handleError(w, "Failed to check existing mood entry", err, http.StatusInternalServerError)
		return
	}

	if entry != nil {
		s.handleError(w, "Mood entry for this date already exists", nil, http.StatusConflict)
		return
	}

	_, err = s.DBOperations.AddMoodEntry(input.UserID, input.Date, input.MoodTypeID, input.Note)
	if err != nil {
		s.handleError(w, "Failed to add mood entry", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, map[string]string{"message": "Mood entry created"}, http.StatusCreated)
}

func (s *Server) handleGetMoodTypes(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood types")

	moodTypes, err := s.DBOperations.GetMoodTypes()
	if err != nil {
		s.handleError(w, "Failed to retrieve mood types", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, moodTypes, http.StatusOK)
}

func (s *Server) handleGetMoods(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting moods")

	input, err := s.parseQueryParams(r)
	if err != nil {
		s.handleError(w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	moods, err := s.DBOperations.GetMoodEntries(*input)
	if err != nil {
		s.handleError(w, "Failed to retrieve moods", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, moods, http.StatusOK)
}

func (s *Server) handleGetMoodSummary(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood summary")

	input, err := s.parseQueryParams(r)
	if err != nil {
		s.handleError(w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	summary, err := s.DBOperations.GetMoodSummary(*input)
	if err != nil {
		s.handleError(w, "Failed to retrieve mood summary", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, summary, http.StatusOK)
}

type updateMoodInput struct {
	ID         int    `json:"id" validate:"required"`
	MoodTypeID int    `json:"moodTypeId" validate:"required"`
	Note       string `json:"note" validate:"required,max=500"`
}

func (s *Server) handleUpdateMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Updating mood entry")
	var input updateMoodInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Struct(input)
	if err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.DBOperations.UpdateMoodEntry(input.ID, input.MoodTypeID, input.Note)
	if err != nil {
		s.handleError(w, "Failed to update mood entry", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, map[string]string{"message": "Mood entry updated"}, http.StatusOK)
}

func (s *Server) handleDeleteMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting mood entry")
	idStr := r.PathValue("id")
	if idStr == "" {
		s.handleError(w, "Missing id parameter", nil, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.handleError(w, "Invalid id parameter", err, http.StatusBadRequest)
		return
	}

	err = s.DBOperations.DeleteMoodEntry(id)
	if err != nil {
		if err.Error() == "no rows deleted" {
			s.handleError(w, "Mood entry not found", err, http.StatusNotFound)
			return
		}
		s.handleError(w, "Failed to delete mood entry", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, map[string]string{"message": "Mood entry deleted"}, http.StatusOK)
}

func (s *Server) handleGetMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood entry by ID")
	idStr := r.PathValue("id")
	if idStr == "" {
		s.handleError(w, "Missing id parameter", nil, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		s.handleError(w, "Invalid id parameter", err, http.StatusBadRequest)
		return
	}

	entry, err := s.DBOperations.GetMoodEntryByID(id)
	if err != nil {
		if err.Error() == "mood entry not found" {
			s.handleError(w, "Mood entry not found", err, http.StatusNotFound)
			return
		}
		s.handleError(w, "Failed to retrieve mood entry", err, http.StatusInternalServerError)
		return
	}

	s.writeJSON(w, entry, http.StatusOK)
}

// Helper function to parse query parameters for mood retrieval (get moods and get summary)
func (s *Server) parseQueryParams(r *http.Request) (*repository.GetInput, error) {
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

	return &repository.GetInput{
		UserID:    userID,
		StartDate: from,
		EndDate:   to,
	}, nil
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
