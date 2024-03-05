package api

import (
	respon "URLShortenePetPrpoject/internal/lib/api/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func GetRedirect(url string) (string, respon.Response, error) {
	const fn = "api.GetRedirect"
	var res respon.Response

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", res, fmt.Errorf("%s: %w", fn, err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", res, fmt.Errorf("%s: %w", fn, err)
	}

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(bodyBytes, &res)
		if err != nil {
			return "", res, fmt.Errorf("%s: %w", fn, err)
		}
	} else if resp.StatusCode != http.StatusFound {
		return "", res, fmt.Errorf("%s, %s: %d", fn, ErrInvalidStatusCode, resp.StatusCode)
	}

	defer func() { _ = resp.Body.Close() }()

	return resp.Header.Get("Location"), res, nil
}
