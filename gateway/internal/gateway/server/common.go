package server

import (
	"io"
	"net/http"
)

// Helper function to handle errors
func (s *Server) handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	if err != nil {
		s.Logger.Error.Printf("%s: %v", message, err)
	} else {
		s.Logger.Error.Println(message)
	}
	http.Error(w, message, statusCode)
}

// Helper function to forward the response
func (s *Server) forwardResponse(w http.ResponseWriter, resp *http.Response) {
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
