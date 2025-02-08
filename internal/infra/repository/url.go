package repository

import (
	"context"
	"github.com/ramk42/mini-url/internal/model"
)

type URL struct {
}

func (U *URL) Save(ctx context.Context, expirationDays int) error {
	//TODO implement me
	panic("implement me")
}

func (U *URL) Get(ctx context.Context, short string) (*model.URL, error) {
	//TODO implement me
	panic("implement me")
}
