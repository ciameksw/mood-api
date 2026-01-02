package externalquotes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/ciameksw/mood-api/quote/internal/quote/config"
)

type ExternalQuotesService struct {
	APIURL string
}

func NewExternalQuotesService(cfg *config.Config) *ExternalQuotesService {
	return &ExternalQuotesService{
		APIURL: cfg.ExternalQuoteURL,
	}
}

type Quote struct {
	Content string   `json:"q"`
	Author  string   `json:"a"`
	Tags    []string `json:"c"`
}

type QuoteResponse struct {
	Quote       string `json:"quote"`
	Author      string `json:"author"`
	Attribution string `json:"attribution"`
}

// GetTodayQuote fetches the daily quote from external API
func (s *ExternalQuotesService) GetTodayQuote() (*QuoteResponse, error) {
	resp, err := http.Get(s.APIURL + "/api/today")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch quote: status " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	quotes := make([]Quote, 0)
	if err := json.Unmarshal(body, &quotes); err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, errors.New("no quotes returned from API")
	}

	quote := quotes[0]
	return &QuoteResponse{
		Quote:       quote.Content,
		Author:      quote.Author,
		Attribution: "Quotes provided by https://zenquotes.io/",
	}, nil
}
