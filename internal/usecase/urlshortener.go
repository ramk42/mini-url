package usecase

import (
	"context"
	"github.com/ramk42/mini-url/internal/model"
	"github.com/ramk42/mini-url/internal/port"
)

type URLShortener struct {
	urlRepo port.URLRepository
}

func NewURLShortener(urlRepo port.URLRepository) *URLShortener {
	return &URLShortener{urlRepo: urlRepo}
}

func (u *URLShortener) ShortenURL(ctx context.Context, url *model.URL) (string, error) {
	return "", nil
}

func (u *URLShortener) Resolve(ctx context.Context, slug string) (*model.URL, error) {
	return nil, nil
}
