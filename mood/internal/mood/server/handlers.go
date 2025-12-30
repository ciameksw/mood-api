package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ciameksw/mood-api/pkg/httputil"
	"github.com/ciameksw/mood-api/pkg/queryutil"
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
		httputil.HandleError(*s.Logger, w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Struct(input)
	if err != nil {
		httputil.HandleError(*s.Logger, w, err.Error(), err, http.StatusBadRequest)
		return
	}

	entry, err := s.DBOperations.GetMoodEntryByDateAndUser(input.UserID, input.Date)
	if err != nil && err.Error() != "user not found" {
		httputil.HandleError(*s.Logger, w, "Failed to check existing mood entry", err, http.StatusInternalServerError)
		return
	}

	if entry != nil {
		httputil.HandleError(*s.Logger, w, "Mood entry for this date already exists", nil, http.StatusConflict)
		return
	}

	_, err = s.DBOperations.AddMoodEntry(input.UserID, input.Date, input.MoodTypeID, input.Note)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to add mood entry", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteSuccessMessage(*s.Logger, w, "Mood entry created", http.StatusCreated)
}

func (s *Server) handleGetMoodTypes(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood types")

	moodTypes, err := s.DBOperations.GetMoodTypes()
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to retrieve mood types", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteData(*s.Logger, w, moodTypes, http.StatusOK)
}

func (s *Server) handleGetMoods(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting moods")

	input, err := queryutil.ParseTimeframeWithUserIDParams(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	moods, err := s.DBOperations.GetMoodEntries(*input)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to retrieve moods", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteData(*s.Logger, w, moods, http.StatusOK)
}

func (s *Server) handleGetMoodSummary(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood summary")

	input, err := queryutil.ParseTimeframeWithUserIDParams(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	summary, err := s.DBOperations.GetMoodSummary(*input)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to retrieve mood summary", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteData(*s.Logger, w, summary, http.StatusOK)
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
		httputil.HandleError(*s.Logger, w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Struct(input)
	if err != nil {
		httputil.HandleError(*s.Logger, w, err.Error(), err, http.StatusBadRequest)
		return
	}

	err = s.DBOperations.UpdateMoodEntry(input.ID, input.MoodTypeID, input.Note)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to update mood entry", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteSuccessMessage(*s.Logger, w, "Mood entry updated", http.StatusOK)
}

func (s *Server) handleDeleteMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting mood entry")
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

	err = s.DBOperations.DeleteMoodEntry(id)
	if err != nil {
		if err.Error() == "no rows deleted" {
			httputil.HandleError(*s.Logger, w, "Mood entry not found", err, http.StatusNotFound)
			return
		}
		httputil.HandleError(*s.Logger, w, "Failed to delete mood entry", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteSuccessMessage(*s.Logger, w, "Mood entry deleted", http.StatusOK)
}

func (s *Server) handleGetMood(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood entry by ID")
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

	entry, err := s.DBOperations.GetMoodEntryByID(id)
	if err != nil {
		if err.Error() == "mood entry not found" {
			httputil.HandleError(*s.Logger, w, "Mood entry not found", err, http.StatusNotFound)
			return
		}
		httputil.HandleError(*s.Logger, w, "Failed to retrieve mood entry", err, http.StatusInternalServerError)
		return
	}

	httputil.WriteData(*s.Logger, w, entry, http.StatusOK)
}
