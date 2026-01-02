package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ciameksw/mood-api/pkg/httputil"
)

type contextKey string

const userIDContextKey contextKey = "userID"

// authMiddleware checks for valid authorization token
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httputil.HandleError(*s.Logger, w, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}

		resp, err := s.AuthService.Authorize(authHeader)
		if err != nil {
			httputil.HandleError(*s.Logger, w, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			httputil.HandleError(*s.Logger, w, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}

		// Parse userId from auth service response
		var body struct {
			UserID int `json:"userId"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			httputil.HandleError(*s.Logger, w, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}
		if body.UserID == 0 {
			httputil.HandleError(*s.Logger, w, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}

		// Attach user id to context and proceed
		ctx := context.WithValue(r.Context(), userIDContextKey, body.UserID)
		next(w, r.WithContext(ctx))
	}
}

// getUserIDFromContext retrieves the authenticated user id set by authMiddleware
func getUserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(userIDContextKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int)
	return id, ok
}
