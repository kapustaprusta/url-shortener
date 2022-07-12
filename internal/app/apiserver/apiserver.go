package apiserver

import (
	"net/http"

	"github.com/kapustaprusta/url-shortener/internal/app/store/dummystore"
)

func Start(config *Config) error {
	store := dummystore.NewStore()
	srv := newServer(store)

	return http.ListenAndServe(config.BindAddr, srv)
}
