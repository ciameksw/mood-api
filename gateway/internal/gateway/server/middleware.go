package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type contextKey string

const userIDContextKey contextKey = "userID"

// authMiddleware checks for valid authorization token
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := s.AuthService.Authorize(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse user_id from auth service response
		var body struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		if body.UserID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
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
