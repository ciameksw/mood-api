package httpclient

import (
	"io"
	"net/http"
)

type RequestParams struct {
	URL           string
	Method        string
	Body          io.Reader
	ContentType   *string
	Authorization *string
}

func SendRequest(params RequestParams) (*http.Response, error) {
	req, err := http.NewRequest(params.Method, params.URL, params.Body)
	if err != nil {
		return nil, err
	}

	if params.ContentType != nil {
		req.Header.Set("Content-Type", *params.ContentType)
	}
	if params.Authorization != nil {
		req.Header.Set("Authorization", *params.Authorization)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
