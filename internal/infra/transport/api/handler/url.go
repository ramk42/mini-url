// Package handler implements the http handlers for the URL shortener service.
package handler

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/ramk42/mini-url/internal/apperr"
	"github.com/ramk42/mini-url/internal/infra/transport/api/httprenderer"
	"github.com/ramk42/mini-url/internal/infra/transport/api/request"
	"github.com/ramk42/mini-url/internal/model"
	"github.com/ramk42/mini-url/internal/port"
	"github.com/ramk42/mini-url/pkg/logger"
	"github.com/ramk42/mini-url/pkg/url/slug"
	"net/http"
	"time"
)

type URL struct {
	shortenerUsecase port.URLShortener
}

func NewURL(shortenerUsecase port.URLShortener) *URL {
	return &URL{shortenerUsecase: shortenerUsecase}
}

func (u *URL) Shorten(w http.ResponseWriter, r *http.Request) {
	req := request.URLShorten{}
	if err := req.Bind(r); err != nil {
		render.Render(w, r, httprenderer.ErrInvalidRequest(err))
		return
	}
	shortenURL, err := u.shortenerUsecase.ShortenURL(r.Context(), model.URL{
		Original:  req.LongURL,
		CreatedAt: time.Now(),
	}, req.ExpirationDays)

	if err != nil {
		render.Render(w, r, httprenderer.ErrUnprocessableEntity(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"short_url": shortenURL})
}

func (u *URL) Resolve(w http.ResponseWriter, r *http.Request) {
	log := logger.Instance().With().Str("method", "Resolve").Logger()
	cleanedSlug := slug.Clean(r.URL.RequestURI())
	err := slug.Validate(cleanedSlug)
	if err != nil {
		log.Error().Err(err).Msg("invalid slug")
		render.Render(w, r, httprenderer.ErrNotFoundRequest(apperr.ErrURLNotFound)) // we don't expose the actual error for security reasons
		return
	}
	resolvedURL, err := u.shortenerUsecase.Resolve(r.Context(), cleanedSlug)
	if err != nil {
		if errors.Is(apperr.ErrURLNotFound, err) {
			render.Render(w, r, httprenderer.ErrNotFoundRequest(err))
			return
		}
		render.Render(w, r, httprenderer.ErrUnprocessableEntity(err))
		return
	}
	http.Redirect(w, r, resolvedURL, http.StatusMovedPermanently)
}
