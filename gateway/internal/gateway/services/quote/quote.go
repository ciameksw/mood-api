package quote

import (
	"net/http"

	"github.com/ciameksw/mood-api/gateway/internal/gateway/config"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/httpclient"
)

type QuoteService struct {
	QuoteURL string
}

func NewQuoteService(cfg *config.Config) *QuoteService {
	return &QuoteService{
		QuoteURL: cfg.QuoteURL,
	}
}

func (qs *QuoteService) GetTodayQuote(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         qs.QuoteURL + "/quote/today",
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
