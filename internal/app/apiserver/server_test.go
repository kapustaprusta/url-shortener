package apiserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kapustaprusta/url-shortener/internal/app/store"
	"github.com/kapustaprusta/url-shortener/internal/app/store/dummystore"
	"github.com/stretchr/testify/require"
)

func TestServerHandleShortenRequest(t *testing.T) {
	type fields struct {
		store store.Store
	}

	type result struct {
		statusCode int
		body       string
	}

	tests := []struct {
		name         string
		url          string
		urlToShorten string
		fields       fields
		result       result
	}{
		{
			name:         "success",
			url:          "/",
			urlToShorten: "http://otl2hdnd.ru/ru88fj",
			fields: fields{
				store: dummystore.NewStore(),
			},
			result: result{
				statusCode: http.StatusCreated,
				body:       "http://localhost:8080/?id=0",
			},
		},
		{
			name:         "empty body",
			url:          "/",
			urlToShorten: "",
			fields: fields{
				store: dummystore.NewStore(),
			},
			result: result{
				statusCode: http.StatusBadRequest,
				body:       "URL to shorten is not valid\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				request := httptest.NewRequest(
					http.MethodPost,
					tt.url,
					strings.NewReader(tt.urlToShorten),
				)
				response := httptest.NewRecorder()
				s := newServer(tt.fields.store)
				s.ServeHTTP(response, request)

				result := response.Result()
				resultBody, err := ioutil.ReadAll(result.Body)
				result.Body.Close()

				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.result.statusCode, result.StatusCode)
				require.Equal(t, tt.result.body, string(resultBody))
			},
		)
	}
}

func TestServerHandleGetOriginalURL(t *testing.T) {
	type fields struct {
		store store.Store
	}

	type result struct {
		statusCode  int
		body        string
		originalURL string
	}

	tests := []struct {
		name         string
		urlPost      string
		urlGet       string
		urlToShorten string
		fields       fields
		result       result
	}{
		{
			name:         "success",
			urlPost:      "/",
			urlGet:       "/?id=0",
			urlToShorten: "http://otl2hdnd.ru/ru88fj",
			fields: fields{
				store: dummystore.NewStore(),
			},
			result: result{
				statusCode:  http.StatusTemporaryRedirect,
				originalURL: "http://otl2hdnd.ru/ru88fj",
			},
		},
		{
			name:         "empty query",
			urlPost:      "/",
			urlGet:       "/",
			urlToShorten: "",
			fields: fields{
				store: dummystore.NewStore(),
			},
			result: result{
				statusCode: http.StatusBadRequest,
				body:       "The query parameter is missing\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name,
			func(t *testing.T) {
				request := httptest.NewRequest(
					http.MethodPost,
					tt.urlPost,
					strings.NewReader(tt.urlToShorten),
				)
				response := httptest.NewRecorder()
				s := newServer(tt.fields.store)
				s.ServeHTTP(response, request)

				request = httptest.NewRequest(
					http.MethodGet,
					tt.urlGet,
					nil,
				)
				response = httptest.NewRecorder()
				s.ServeHTTP(response, request)

				result := response.Result()
				resultBody, err := ioutil.ReadAll(result.Body)
				result.Body.Close()

				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, tt.result.body, string(resultBody))
				require.Equal(t, tt.result.statusCode, result.StatusCode)
				require.Equal(t, tt.result.originalURL, result.Header.Get("Location"))
			},
		)
	}
}
