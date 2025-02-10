// Package usecase contains the business logic for the URLShortener
package usecase

import (
	"context"
	"errors"
	"github.com/avast/retry-go"
	"github.com/ramk42/mini-url/internal/apperr"
	"github.com/ramk42/mini-url/internal/model"
	"github.com/ramk42/mini-url/internal/port"
	"github.com/ramk42/mini-url/pkg/logger"
	"github.com/ramk42/mini-url/pkg/url/normalizer"
	"github.com/ramk42/mini-url/pkg/url/slug"
	"time"
)

type URLShortener struct {
	urlRepo port.URLRepository
	baseURL string
}

func NewURLShortener(urlRepo port.URLRepository, baseURL string) *URLShortener {
	return &URLShortener{urlRepo: urlRepo, baseURL: baseURL}
}

func (u *URLShortener) ShortenURL(ctx context.Context, url model.URL, expirationDays int) (string, error) {
	log := logger.Instance().With().Str("method", "ShortenURL").Logger()
	normalizedURL, err := normalizer.NormalizeURL(url.Original)
	if err != nil {
		log.Error().Str("url", url.Original).Err(err).Msg("failure to normalize url")
	}
	url.Original = normalizedURL
	if expirationDays > 0 {
		t := time.Now().AddDate(0, 0, expirationDays)
		url.ExpiresAt = &t
	}
	var persistedURL model.URL
	err = retry.Do(
		func() error {
			url.Slug = slug.Generate(model.Sluglength)
			persistedURL, err = u.urlRepo.Save(ctx, url)
			if err != nil {
				if errors.Is(err, apperr.ErrURLConflict) {
					log.Warn().Msg("slug collision detected! retrying...")
					return err
				}
			}
			return err
		},
		retry.Attempts(3),
		retry.Delay(20*time.Millisecond),
	)
	if err != nil {
		log.Error().Err(err).Msg("failure to generate unique slug")
		return "", err
	}

	return u.baseURL + "/" + persistedURL.Slug, nil
}

func (u *URLShortener) Resolve(ctx context.Context, slug string) (string, error) {
	log := logger.Instance().With().Str("method", "Resolve").Logger()
	persistedURL, err := u.urlRepo.UpdateClicks(ctx, slug)
	if err != nil {
		log.Error().Err(err).Str("slug", slug).Msg("failure to get url")
		return "", apperr.ErrURLNotFound
	}
	return persistedURL.Original, nil
}
