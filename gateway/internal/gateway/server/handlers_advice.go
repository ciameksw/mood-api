package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ciameksw/mood-api/pkg/httputil"
)

type selectAdviceInputEntry struct {
	MoodTypeID int     `json:"moodTypeId" validate:"required"`
	Count      int     `json:"count" validate:"required,min=1"`
	Percentage float64 `json:"percentage" validate:"required"`
}

func (s *Server) handleGetAdvice(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting advice")

	from, to, err := s.parseQueryParams(r)
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

	if resp.StatusCode != http.StatusOK {
		s.forwardResponse(w, resp)
		return
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to read mood summary", err, http.StatusInternalServerError)
		return
	}

	var entries []selectAdviceInputEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to parse mood summary", err, http.StatusInternalServerError)
		return
	}

	// Validate entries according to tags (optional but useful)
	for _, e := range entries {
		if err := s.Validator.Struct(e); err != nil {
			httputil.HandleError(*s.Logger, w, "Invalid mood summary entry", err, http.StatusBadRequest)
			return
		}
	}

	// Send the parsed summary entries to the advice service's select endpoint
	bodyBytes, err := json.Marshal(entries)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to marshal request body", err, http.StatusInternalServerError)
		return
	}

	updateResp, err := s.AdviceService.Select(bodyBytes)
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to send update request to advice service", err, http.StatusInternalServerError)
		return
	}
	s.forwardResponse(w, updateResp)
}
