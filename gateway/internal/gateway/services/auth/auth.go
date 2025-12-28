package auth

import (
	"net/http"

	"github.com/ciameksw/mood-api/gateway/internal/gateway/config"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/httpclient"
)

type AuthService struct {
	AuthURL string
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{
		AuthURL: cfg.AuthURL,
	}
}

func (as *AuthService) Register(r *http.Request) (*http.Response, error) {
	return as.commonServiceFunc("/auth/register", r)
}

func (as *AuthService) Login(r *http.Request) (*http.Response, error) {
	return as.commonServiceFunc("/auth/login", r)
}

func (as *AuthService) Authorize(authHeader string) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:           as.AuthURL + "/auth/authorize",
		Method:        http.MethodGet,
		Authorization: &authHeader,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (as *AuthService) GetLoggedUser(r *http.Request) (*http.Response, error) {
	return as.commonServiceFunc("/auth/user", r)
}

func (as *AuthService) UpdateLoggedUser(r *http.Request) (*http.Response, error) {
	return as.commonServiceFunc("/auth/user", r)
}

func (as *AuthService) DeleteLoggedUser(r *http.Request) (*http.Response, error) {
	return as.commonServiceFunc("/auth/user", r)
}

func (as *AuthService) commonServiceFunc(url string, r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         as.AuthURL + url,
		Method:      r.Method,
		Body:        r.Body,
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
