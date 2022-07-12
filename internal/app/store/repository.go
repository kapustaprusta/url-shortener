package store

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
)

type URLRepository interface {
	Add(*model.URL) error
	FindByID(int) (*model.URL, error)
}
