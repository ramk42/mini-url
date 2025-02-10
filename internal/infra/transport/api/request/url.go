// Package request bind request and validate the URLShorten struct
package request

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
)

type URLShorten struct {
	LongURL        string `json:"long_url" validate:"required"`
	ExpirationDays int    `json:"expiration_days" validate:"omitempty,gte=0"`
}

func (u *URLShorten) Bind(r *http.Request) error {
	// Parse the JSON body
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		return err
	}

	// Validate the struct fields
	validate := validator.New()
	if err := validate.Struct(u); err != nil {
		return err
	}

	// Validate the URL format
	if _, err := url.ParseRequestURI(u.LongURL); err != nil {
		return err
	}

	parsedURL, err := url.Parse(u.LongURL)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return errors.New("invalid Scheme or Host")
	}

	return nil
}
