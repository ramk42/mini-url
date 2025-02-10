// Package repository provides the database operations for the URL shortener service.
package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ramk42/mini-url/internal/apperr"
	"github.com/ramk42/mini-url/internal/model"
	"github.com/ramk42/mini-url/pkg/logger"
	"time"
)

type URL struct {
	db *sql.DB
}

func NewURL(db *sql.DB) *URL {
	return &URL{db: db}
}

type URLPersisted struct {
	ID            int
	Original      string
	NormalizedURL string
	Slug          string
	ExpiresAt     sql.NullTime
	Clicks        int
}

func (u *URL) Save(ctx context.Context, url model.URL) (model.URL, error) {
	log := logger.Instance().With().Str("method", "Save").Logger()
	query := `
		INSERT INTO urls (original_url, slug, expires_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (original_url) 
		DO UPDATE SET original_url = EXCLUDED.original_url
		RETURNING slug;
	`

	var urlPersisted URLPersisted
	err := u.db.QueryRowContext(
		ctx,
		query,
		url.Original,
		url.Slug,
		url.ExpiresAt,
	).Scan(
		&urlPersisted.Slug,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return model.URL{}, apperr.ErrURLConflict
		}
		log.Error().Err(err).Msg("failure to save url")
		return model.URL{}, apperr.ErrURLInsertFailed
	}
	return model.URL{
		Slug: urlPersisted.Slug,
	}, nil
}

func (u *URL) Get(ctx context.Context, slug string) (model.URL, error) {
	log := logger.Instance().With().Str("method", "Get").Logger()
	query := `SELECT original_url FROM urls WHERE slug = $1 AND (expires_at IS NULL OR expires_at > $2);`
	var urlPersisted URLPersisted
	err := u.db.QueryRowContext(ctx, query, slug, time.Now()).Scan(
		&urlPersisted.Original,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("url not found")
			return model.URL{}, apperr.ErrURLNotFound
		}
		log.Error().Err(err).Msg("failure to get url")
		return model.URL{}, apperr.ErrURLFetchFailed
	}

	return model.URL{
		Original: urlPersisted.Original,
	}, nil
}

func (u *URL) UpdateClicks(ctx context.Context, slug string) (model.URL, error) {
	log := logger.Instance().With().Str("method", "UpdateClicks").Logger()
	query := `
		UPDATE urls
		SET clicks = clicks + 1
		WHERE slug = $1 AND (expires_at IS NULL OR expires_at > $2)
		RETURNING original_url;
	`
	var urlPersisted URLPersisted

	err := u.db.QueryRowContext(ctx, query, slug, time.Now()).Scan(&urlPersisted.Original)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("url not found")
			return model.URL{}, apperr.ErrURLNotFound
		}
		log.Error().Err(err).Msg("failure to update clicks")
		return model.URL{}, apperr.ErrURLUpdateClicksFailed
	}

	return model.URL{
		Original: urlPersisted.Original,
	}, nil
}
