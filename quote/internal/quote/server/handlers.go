package server

import (
	"context"
	"net/http"

	"github.com/ciameksw/mood-api/pkg/httputil"
)

func (s *Server) handleGetTodayQuote(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting today's quote")

	ctx := context.Background()

	cachedQuote, err := s.RedisCache.GetTodayQuote(ctx)
	if err == nil && cachedQuote != nil {
		s.Logger.Info.Println("Quote found in cache")
		httputil.WriteData(*s.Logger, w, cachedQuote, http.StatusOK)
		return
	}

	resp, err := s.ExternalQuotesService.GetTodayQuote()
	if err != nil {
		httputil.HandleError(*s.Logger, w, "Failed to get today's quote", err, http.StatusInternalServerError)
		return
	}

	if err := s.RedisCache.SetTodayQuote(ctx, resp); err != nil {
		s.Logger.Error.Printf("Failed to cache quote: %v", err)
	}

	httputil.WriteData(*s.Logger, w, resp, http.StatusOK)
}
