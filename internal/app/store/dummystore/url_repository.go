package dummystore

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
	"github.com/kapustaprusta/url-shortener/internal/app/store"
)

type URLRepository struct {
	store *Store
	urls  map[int]*model.URL
}

func (r *URLRepository) Add(u *model.URL) error {
	u.ID = len(r.urls)
	r.urls[u.ID] = u

	return nil
}

func (r *URLRepository) FindById(shortenedURLId int) (*model.URL, error) {
	u, isOk := r.urls[shortenedURLId]
	if !isOk {
		return nil, store.ErrRecordNotFound
	}

	return u, nil
}
