package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/ciameksw/mood-api/pkg/httputil"
	"github.com/ciameksw/mood-api/pkg/queryutil"
)

type addMoodInput struct {
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

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	body := map[string]interface{}{
		"userId":     userID,
		"moodTypeId": input.MoodTypeID,
		"note":       input.Note,
		"date":       input.Date,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to marshal request body", err, http.StatusInternalServerError)
		return
	}

	resp, err := s.MoodService.Add(bodyBytes)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleGetMoodTypes(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Get mood types")

	resp, err := s.MoodService.GetTypes(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleGetMoodSummary(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting mood summary")

	from, to, err := queryutil.ParseTimeframeParams(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.MoodService.GetSummary(from, to, userID)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleGetMoods(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting moods")

	from, to, err := queryutil.ParseTimeframeParams(r)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Invalid query parameters", err, http.StatusBadRequest)
		return
	}

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.MoodService.GetMoods(from, to, userID)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
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

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.MoodService.GetMood(id)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.forwardResponse(w, resp)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to read mood entry", err, http.StatusInternalServerError)
		return
	}

	var entry struct {
		UserID int `json:"UserID"`
	}
	if err := json.Unmarshal(bodyBytes, &entry); err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to decode mood entry", err, http.StatusInternalServerError)
		return
	}

	if entry.UserID != userID {
		httputil.HandleError(*s.Logger, w, "Forbidden: mood entry does not belong to user", nil, http.StatusForbidden)
		return
	}

	httputil.WriteData(*s.Logger, w, entry, http.StatusOK)
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

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.MoodService.GetMood(id)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.forwardResponse(w, resp)
		return
	}

	var entry struct {
		UserID int `json:"UserID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		resp.Body.Close()
		httputil.HandleError(*s.Logger, w, "Failed to decode mood entry", err, http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	if entry.UserID != userID {
		httputil.HandleError(*s.Logger, w, "Forbidden: mood entry does not belong to user", nil, http.StatusForbidden)
		return
	}

	delResp, err := s.MoodService.DeleteMood(id)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send delete request to mood service", err, http.StatusInternalServerError)
		return
	}
	s.forwardResponse(w, delResp)
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

	userID, ok := getUserIDFromContext(r.Context())
	if !ok {
		httputil.HandleError(*s.Logger, w, "Failed to get user ID from context", nil, http.StatusUnauthorized)
		return
	}

	resp, err := s.MoodService.GetMood(input.ID)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send request to mood service", err, http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.forwardResponse(w, resp)
		return
	}

	var entry struct {
		UserID int `json:"UserID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		resp.Body.Close()
		httputil.HandleError(*s.Logger, w, "Failed to decode mood entry", err, http.StatusInternalServerError)
		return
	}
	resp.Body.Close()

	if entry.UserID != userID {
		httputil.HandleError(*s.Logger, w, "Forbidden: mood entry does not belong to user", nil, http.StatusForbidden)
		return
	}

	body := map[string]interface{}{
		"id":         input.ID,
		"moodTypeId": input.MoodTypeID,
		"note":       input.Note,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to marshal request body", err, http.StatusInternalServerError)
		return
	}

	updateResp, err := s.MoodService.Update(bodyBytes)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send update request to mood service", err, http.StatusInternalServerError)
		return
	}
	s.forwardResponse(w, updateResp)
}
