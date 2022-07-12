package store

import (
	"github.com/kapustaprusta/url-shortener/internal/app/model"
)

type UrlRepository interface {
	Add(*model.Url) error
	FindByShortenedUrl(string) (*model.Url, error)
}
