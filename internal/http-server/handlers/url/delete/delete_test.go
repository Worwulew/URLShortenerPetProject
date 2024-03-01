package delete_test

import (
	deleter "URLShortenePetPrpoject/internal/http-server/handlers/url/delete"
	"URLShortenePetPrpoject/internal/http-server/handlers/url/delete/mocks"
	"URLShortenePetPrpoject/internal/lib/logger/handlers/slogdiscard"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
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
			name:      "DeleteURL error",
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

			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleterMock.On("DeleteURL", tc.alias).
					Return(tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", deleter.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			client := &http.Client{}

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/"+tc.alias, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
