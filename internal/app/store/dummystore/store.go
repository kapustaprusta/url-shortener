package dummystore

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

type Store struct {
	urlRepository *UrlRepository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) Url() store.UrlRepository {
	if s.urlRepository != nil {
		return s.urlRepository
	}

	s.urlRepository = &UrlRepository{
		store: s,
		urls:  make(map[string]*model.Url),
	}

	return s.urlRepository
}
