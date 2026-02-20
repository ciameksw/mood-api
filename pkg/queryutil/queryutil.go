package queryutil

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

func ParseTimeframeParams(r *http.Request) (string, string, error) {
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	if from == "" || to == "" {
		return "", "", errors.New("from and to parameters are required")
	}

	// Validate date format YYYY-MM-DD
	const dateFormat = "2006-01-02"
	fromDate, err := time.Parse(dateFormat, from)
	if err != nil {
		return "", "", errors.New("from date must be in YYYY-MM-DD format")
	}
	toDate, err := time.Parse(dateFormat, to)
	if err != nil {
		return "", "", errors.New("to date must be in YYYY-MM-DD format")
	}

	// Validate that from is before to
	if fromDate.After(toDate) {
		return "", "", errors.New("from date must be before or equal to to date")
	}

	return from, to, nil
}

type GetParams struct {
	UserID    int
	StartDate string
	EndDate   string
}

func ParseTimeframeWithUserIDParams(r *http.Request) (*GetParams, error) {
	userID, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		return nil, err
	}

	from, to, err := ParseTimeframeParams(r)
	if err != nil {
		return nil, err
	}

	return &GetParams{
		UserID:    userID,
		StartDate: from,
		EndDate:   to,
	}, nil
}
