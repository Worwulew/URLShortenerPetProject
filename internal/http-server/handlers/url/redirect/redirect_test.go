package redirect_test

import (
	"URLShortenePetPrpoject/internal/http-server/handlers/url/redirect"
	"URLShortenePetPrpoject/internal/http-server/handlers/url/redirect/mocks"
	"URLShortenePetPrpoject/internal/lib/api"
	"URLShortenePetPrpoject/internal/lib/logger/handlers/slogdiscard"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestGetHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:      "GetURL error",
			alias:     "test_alias",
			url:       "https://google.com",
			respError: "internal error",
			mockError: errors.New("unexpected error"),
		},
	}
	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, resp, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)
			if redirectedToURL == "" {
				require.Equal(t, tc.respError, resp.Error)
				return
			}
			require.Equal(t, tc.url, redirectedToURL)
		})
	}
}
