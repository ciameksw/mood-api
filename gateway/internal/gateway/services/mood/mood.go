package mood

import (
	"bytes"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ciameksw/mood-api/gateway/internal/gateway/config"
	"github.com/ciameksw/mood-api/gateway/internal/gateway/httpclient"
)

type MoodService struct {
	MoodURL string
}

func NewMoodService(cfg *config.Config) *MoodService {
	return &MoodService{
		MoodURL: cfg.MoodURL,
	}
}

func (ms *MoodService) Add(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         ms.MoodURL + "/mood",
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

func (ms *MoodService) GetTypes(r *http.Request) (*http.Response, error) {
	ct := r.Header.Get("Content-Type")
	params := httpclient.RequestParams{
		URL:         ms.MoodURL + "/mood/types",
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

func (ms *MoodService) GetSummary(from, to string, userID int) (*http.Response, error) {
	q := url.Values{}
	q.Set("from", from)
	q.Set("to", to)
	q.Set("userId", strconv.Itoa(userID))

	params := httpclient.RequestParams{
		URL:    ms.MoodURL + "/mood/summary?" + q.Encode(),
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ms *MoodService) GetMoods(from, to string, userID int) (*http.Response, error) {
	q := url.Values{}
	q.Set("from", from)
	q.Set("to", to)
	q.Set("userId", strconv.Itoa(userID))

	params := httpclient.RequestParams{
		URL:    ms.MoodURL + "/mood?" + q.Encode(),
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ms *MoodService) DeleteMood(moodID int) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    ms.MoodURL + "/mood/" + strconv.Itoa(moodID),
		Method: http.MethodDelete,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ms *MoodService) GetMood(moodID int) (*http.Response, error) {
	params := httpclient.RequestParams{
		URL:    ms.MoodURL + "/mood/" + strconv.Itoa(moodID),
		Method: http.MethodGet,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (ms *MoodService) Update(body []byte) (*http.Response, error) {
	ct := "application/json"
	params := httpclient.RequestParams{
		URL:         ms.MoodURL + "/mood",
		Method:      http.MethodPut,
		Body:        bytes.NewBuffer(body),
		ContentType: &ct,
	}
	resp, err := httpclient.SendRequest(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
