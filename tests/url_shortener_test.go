package tests

import (
	"URLShortenePetPrpoject/internal/http-server/handlers/url/save"
	"URLShortenePetPrpoject/internal/lib/api"
	"URLShortenePetPrpoject/internal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

const host = "localhost:8082"

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(11),
		}).
		WithBasicAuth("myUser", "1924").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirectDelete(t *testing.T) {
	testCases := []struct {
		name  string
		alias string
		url   string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid url",
		},
		{
			name:  "Empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// save

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myUser", "1924").
				Expect().Status(200).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			alias := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			// redirect

			testRedirect(t, alias, tc.url)

			// delete

			respDel := e.DELETE("/url/"+alias).
				WithBasicAuth("myUser", "1924").
				Expect().Status(200).
				JSON().Object()

			respDel.Value("status").String().IsEqual("Ok")

			// redirect not found

			testRedirectNotFound(t, alias)
		})
	}
}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "url/" + alias,
	}

	redirectedURL, _, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedURL)
}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   "url/" + alias,
	}

	_, resp, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, resp.Error, "not found")
}
