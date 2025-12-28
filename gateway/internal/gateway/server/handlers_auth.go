package server

import "net/http"

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Register user")

	resp, err := s.AuthService.Register(r)
	if err != nil {
		s.handleError(w, "Failed to send request to auth service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Login user")

	resp, err := s.AuthService.Login(r)
	if err != nil {
		s.handleError(w, "Failed to send request to auth service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Get logged user")

	resp, err := s.AuthService.GetLoggedUser(r)
	if err != nil {
		s.handleError(w, "Failed to send request to auth service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Update logged user")

	resp, err := s.AuthService.UpdateLoggedUser(r)
	if err != nil {
		s.handleError(w, "Failed to send request to auth service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	s.Logger.Info.Println("Delete logged user")

	resp, err := s.AuthService.DeleteLoggedUser(r)
	if err != nil {
		s.handleError(w, "Failed to send request to auth service", err, http.StatusInternalServerError)
		return
	}

	s.forwardResponse(w, resp)
}
