package advice

import (
	"bytes"
	"net/http"

	"github.com/ciameksw/mood-api/gateway/internal/gateway/config"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/httpclient"
)

type AdviceService struct {
	AdviceURL string
}

func NewAdviceService(cfg *config.Config) *AdviceService {
	return &AdviceService{
		AdviceURL: cfg.AdviceURL,
	}
}

func (as *AdviceService) Select(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         as.AdviceURL + "/advice/select",
		Method:      http.MethodPost,
		Body:        bytes.NewBuffer(body),
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
