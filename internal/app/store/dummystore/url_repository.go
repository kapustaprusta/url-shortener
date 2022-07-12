package dummystore

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

type UrlRepository struct {
	store *Store
	urls  map[string]*model.Url
}

func (r *UrlRepository) Add(u *model.Url) error {
	r.urls[u.ShortenedUrl] = u
	u.ID = len(r.urls)

	return nil
}

func (r *UrlRepository) FindByShortenedUrl(shortenedUrl string) (*model.Url, error) {
	u, isOk := r.urls[shortenedUrl]
	if !isOk {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
