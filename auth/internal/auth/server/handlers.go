package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

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

	existingUser, err := s.DBOperations.GetUserByEmail(r.Context(), input.Email)
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

	userID, err := s.DBOperations.CreateUser(r.Context(), input.UserName, input.Email, hashedPassword)
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

	user, err := s.DBOperations.GetUserByEmail(r.Context(), input.Email)
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
	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.handleError(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"user_id": userID,
	}
	s.writeJSON(w, resp, http.StatusOK)
}

type userResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Getting user profile")

	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.handleError(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	user, err := s.DBOperations.GetUserByID(r.Context(), userID)
	if err != nil {
		s.handleError(w, "Failed to retrieve user", err, http.StatusInternalServerError)
		return
	}

	if user == nil {
		s.handleError(w, "User not found", nil, http.StatusNotFound)
		return
	}

	resp := userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	s.writeJSON(w, resp, http.StatusOK)
}

type updateUserInput struct {
	Username string `json:"username" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Updating user profile")

	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.handleError(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	var input updateUserInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		s.handleError(w, "Invalid request payload", err, http.StatusBadRequest)
		return
	}

	err = s.Validator.Struct(input)
	if err != nil {
		s.handleError(w, err.Error(), err, http.StatusBadRequest)
		return
	}

	// Check if new email already exists
	if input.Email != "" {
		existingUser, err := s.DBOperations.GetUserByEmail(r.Context(), input.Email)
		if err != nil && err.Error() != "user not found" {
			s.handleError(w, "Failed to check existing email", err, http.StatusInternalServerError)
			return
		}
		if existingUser != nil && existingUser.ID != userID {
			s.handleError(w, "Email already in use", nil, http.StatusConflict)
			return
		}
	}

	// Check if new username already exists
	if input.Username != "" {
		existingUser, err := s.DBOperations.GetUserByUsername(r.Context(), input.Username)
		if err != nil && err.Error() != "user not found" {
			s.handleError(w, "Failed to check existing username", err, http.StatusInternalServerError)
			return
		}
		if existingUser != nil && existingUser.ID != userID {
			s.handleError(w, "Username already in use", nil, http.StatusConflict)
			return
		}
	}

	// Hash password if provided
	var hashedPassword *string
	if input.Password != "" {
		hashed, err := token.HashPassword(input.Password)
		if err != nil {
			s.handleError(w, "Failed to hash password", err, http.StatusInternalServerError)
			return
		}
		hashedPassword = &hashed
	}

	err = s.DBOperations.UpdateUser(r.Context(), userID, input.Username, input.Email, hashedPassword)
	if err != nil {
		if err.Error() == "no fields to update" {
			s.handleError(w, "No fields to update", nil, http.StatusBadRequest)
			return
		}
		s.handleError(w, "Failed to update user", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User updated: %d", userID)
	s.writeJSON(w, map[string]string{"message": "User updated successfully"}, http.StatusOK)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Deleting user account")

	userID, err := s.getUserIDFromToken(r)
	if err != nil {
		s.handleError(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	err = s.DBOperations.DeleteUser(r.Context(), userID)
	if err != nil {
		s.handleError(w, "Failed to delete user", err, http.StatusInternalServerError)
		return
	}

	s.Logger.Info.Printf("User deleted: %d", userID)
	s.writeJSON(w, map[string]string{"message": "User deleted successfully"}, http.StatusOK)
}

// Helper function to extract userID from Authorization header
func (s *Server) getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, errors.New("missing Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := token.ValidateJWT(tokenString, s.Config.Salt)
	if err != nil {
		return 0, errors.New("invalid or expired token")
	}

	return claims.UserID, nil
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
