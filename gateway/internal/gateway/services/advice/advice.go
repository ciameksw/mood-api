package advice

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"

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

func (as *AdviceService) SavePeriod(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         as.AdviceURL + "/advice/period/save",
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

func (as *AdviceService) GetByPeriod(from, to string, userID int) (*http.Response, error) {
	q := url.Values{}
	q.Set("from", from)
	q.Set("to", to)
	q.Set("userId", strconv.Itoa(userID))

	params := httpclient.RequestParams{
		URL:    as.AdviceURL + "/advice/period/get?" + q.Encode(),
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
