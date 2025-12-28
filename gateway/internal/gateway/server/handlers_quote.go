package server

import "net/http"

func (s *Server) handleGetTodayQuote(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Get today's quote")

	resp, err := s.QuoteService.GetTodayQuote(r)
	if err != nil {
		s.handleError(w, "Failed to send request to quote service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}
