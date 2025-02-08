package port

import (
	"context"
	"github.com/ramk42/mini-url/internal/model"
)

type URLShortener interface {
	ShortenURL(ctx context.Context, url *model.URL) (string, error)
	Resolve(ctx context.Context, slug string) (*model.URL, error)
}

type URLRepository interface {
	Save(ctx context.Context, expirationDays int) error
	Get(ctx context.Context, shortUrl string) (*model.URL, error)
}
