package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ciameksw/mood-api/auth/internal/auth/token"
)

type registerInput struct {
	UserName string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Registering user")
	var input registerInput

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

	existingUser, err := s.Postgres.GetUserByEmail(r.Context(), input.Email)
	if err != nil && err.Error() != "user not found" {
		s.handleError(w, "Failed to check existing user", err, http.StatusInternalServerError)
		return
	}

	if existingUser != nil {
		s.handleError(w, "User with this email already exists", nil, http.StatusConflict)
		return
	}

	hashedPassword, err := token.HashPassword(input.Password)
	if err != nil {
		s.handleError(w, "Failed to hash password", err, http.StatusInternalServerError)
		return
	}

	userID, err := s.Postgres.CreateUser(r.Context(), input.UserName, input.Email, hashedPassword)
	if err != nil {
		s.handleError(w, "Failed to create user", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User registered successfully id: %d, username: %s, email: %s", userID, input.UserName, input.Email)
	s.writeJSON(w, map[string]string{"message": "User registered successfully"}, http.StatusCreated)
}

type loginInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Logging user")
	var input loginInput

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

	user, err := s.Postgres.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		if err.Error() == "user not found" {
			s.handleError(w, "Invalid email or password", nil, http.StatusUnauthorized)
			return
		}

		s.handleError(w, "Unexpected server error", err, http.StatusInternalServerError)
		return
	}

	if user == nil {
		s.handleError(w, "Invalid email or password", nil, http.StatusUnauthorized)
		return
	}

	match := token.VerifyPassword(input.Password, user.PasswordHash)
	if !match {
		s.handleError(w, "Invalid username or password", err, http.StatusUnauthorized)
		return
	}

	token, err := token.GenerateJWT(user.ID, s.Config.Salt)
	if err != nil {
		s.handleError(w, "Failed to generate JWT", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User logged in: %v", user.Username)
	resp := loginResponse{
		Token: token,
	}
	s.writeJSON(w, resp, http.StatusOK)
}

func (s *Server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		s.handleError(w, "Missing Authorization header", nil, http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := token.ValidateJWT(tokenString, s.Config.Salt)
	if err != nil {
		s.handleError(w, "Invalid or expired token", err, http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"user_id": claims.UserID,
	}
	s.writeJSON(w, resp, http.StatusOK)
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
