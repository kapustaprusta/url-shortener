package apiserver

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

const (
	baseURL = "http://localhost:8080/"
)

type server struct {
	router *mux.Router
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
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
		s.handleGetOriginalURL(w, r)
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
	rawOriginalURL, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(
			w,
			"Something went wrong. Please, try again later",
			http.StatusInternalServerError,
		)

		return
	}
	defer r.Body.Close()

	originalURL := string(rawOriginalURL)
	_, err = url.ParseRequestURI(originalURL)
	if err != nil {
		http.Error(
			w,
			"URL to shorten is not valid",
			http.StatusBadRequest,
		)

		return
	}

	url := &model.URL{
		OriginalURL:  originalURL,
		ShortenedURL: baseURL,
	}

	s.store.URL().Add(url)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(url.ShortenedURL + "?id=" + strconv.Itoa(url.ID)))
}

func (s *server) handleGetOriginalURL(w http.ResponseWriter, r *http.Request) {
	queryShortenedURLId := r.URL.Query().Get("id")
	if queryShortenedURLId == "" {
		http.Error(
			w,
			"The query parameter is missing",
			http.StatusBadRequest,
		)

		return
	}

	shortenedURLId, err := strconv.Atoi(queryShortenedURLId)
	if err != nil {
		http.Error(
			w,
			"Invalid shortened URL ID",
			http.StatusBadRequest,
		)

		return
	}

	url, err := s.store.URL().FindByID(shortenedURLId)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)

		return
	}

	w.Header().Set("Location", url.OriginalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
