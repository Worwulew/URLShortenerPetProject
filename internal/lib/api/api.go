package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func GetRedirect(url string) (string, error) {
	const fn = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s, %s: %d", fn, ErrInvalidStatusCode, resp.StatusCode)
	}

	defer func() { _ = resp.Body.Close() }()

	return resp.Header.Get("Location"), nil
}
