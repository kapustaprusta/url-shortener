package dummystore

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

type URLRepository struct {
	store *Store
	urls  map[string]*model.URL
}

func (r *URLRepository) Add(u *model.URL) error {
	r.urls[u.ShortenedURL] = u
	u.ID = len(r.urls)

	return nil
}

func (r *URLRepository) FindByShortenedURL(shortenedRL string) (*model.URL, error) {
	u, isOk := r.urls[shortenedRL]
	if !isOk {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
