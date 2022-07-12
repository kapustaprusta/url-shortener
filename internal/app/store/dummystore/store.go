package dummystore

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

type Store struct {
	urlRepository *URLRepository
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) URL() store.URLRepository {
	if s.urlRepository != nil {
		return s.urlRepository
	}

	s.urlRepository = &URLRepository{
		store: s,
		urls:  make(map[string]*model.URL),
	}

	return s.urlRepository
}
