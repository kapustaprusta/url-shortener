package store

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
)

type URLRepository interface {
	Add(*model.URL) error
	FindById(int) (*model.URL, error)
}
