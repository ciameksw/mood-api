package server

import "net/http"

func (s *Server) setupAuthRouter(r *http.ServeMux) {
	r.HandleFunc("POST /auth/register", s.handleRegister)                   // Register to the system
	r.HandleFunc("POST /auth/login", s.handleLogin)                         // Login to get auth token
	r.HandleFunc("GET /auth/user", s.authMiddleware(s.handleGetUser))       // Get logged user info
	r.HandleFunc("PUT /auth/user", s.authMiddleware(s.handleUpdateUser))    // Update logged user info
	r.HandleFunc("DELETE /auth/user", s.authMiddleware(s.handleDeleteUser)) // Delete logged user account
}

func (s *Server) setupMoodRouter(r *http.ServeMux) {
	r.HandleFunc("POST /mood", s.authMiddleware(s.handleAddMood))               // Add new mood entry to the logged user
	r.HandleFunc("GET /mood", s.authMiddleware(s.handleGetMoods))               // Get mood entries of the logged user in time range
	r.HandleFunc("GET /mood/types", s.authMiddleware(s.handleGetMoodTypes))     // Get all available mood types
	r.HandleFunc("GET /mood/summary", s.authMiddleware(s.handleGetMoodSummary)) // Get mood summary for the logged user in time range
	r.HandleFunc("GET /mood/{id}", s.authMiddleware(s.handleGetMood))           // Get single mood entry by id
	r.HandleFunc("PUT /mood", s.authMiddleware(s.handleUpdateMood))             // Update a mood entry of the logged user
	r.HandleFunc("DELETE /mood/{id}", s.authMiddleware(s.handleDeleteMood))     // Delete a mood entry of the logged user
}

func (s *Server) setupAdviceRouter(r *http.ServeMux) {
	r.HandleFunc("GET /advice", s.authMiddleware(s.handleGetAdvice)) // Get advice for the logged user
}

func (s *Server) setupQuoteRouter(r *http.ServeMux) {
	r.HandleFunc("GET /quote/today", s.authMiddleware(s.handleGetTodayQuote)) // Get todays quote for the logged user
}
