package port

import (
	"context"
	"github.com/ramk42/mini-url/internal/model"
)

type (
	URLShortener interface {
		ShortenURL(ctx context.Context, url model.URL, expirationInDays int) (string, error)
		Resolve(ctx context.Context, slug string) (string, error)
	}

	URLRepository interface {
		Save(ctx context.Context, url model.URL) (model.URL, error)
		Get(ctx context.Context, slug string) (model.URL, error)
		UpdateClicks(ctx context.Context, slug string) (model.URL, error)
	}
)
