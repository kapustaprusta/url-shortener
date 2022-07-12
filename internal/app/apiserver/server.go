package apiserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

const (
	baseUrl = "www.link.com/"
)

type server struct {
	router *http.ServeMux
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: http.NewServeMux(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.handleRoot)
}

func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetOriginalUrl(w, r)
	case http.MethodPost:
		s.handleShortenRequest(w, r)
	default:
		http.Error(
			w,
			"Only GET and POST methods are allowed",
			http.StatusBadRequest,
		)
	}
}

func (s *server) handleShortenRequest(w http.ResponseWriter, r *http.Request) {
	rawOriginalUrl, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			"Something went wrong. Please, try again later",
			http.StatusInternalServerError,
		)

		return
	}
	defer r.Body.Close()

	originalUrl := string(rawOriginalUrl)
	_, err = url.ParseRequestURI(originalUrl)
	if err != nil {
		http.Error(
			w,
			"URL to shorten is not valid",
			http.StatusBadRequest,
		)

		return
	}

	url := &model.Url{
		OriginalUrl: originalUrl,
		ShortenedUrl: fmt.Sprint(
			baseUrl + originalUrl[0:len(originalUrl)/2],
		),
	}

	s.store.Url().Add(url)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortenedUrl))
}

func (s *server) handleGetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := r.URL.Query().Get("id")
	if shortenedUrl == "" {
		http.Error(
			w,
			"The query parameter is missing",
			http.StatusBadRequest,
		)

		return
	}

	url, err := s.store.Url().FindByShortenedUrl(shortenedUrl)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.WriteHeader(http.StatusTemporaryRedirect)
	w.Header().Set("Location", url.OriginalUrl)
}
