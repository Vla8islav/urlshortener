package storage

import (
	"context"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
)

type URLPair struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
	Deleted  bool   `json:"deleted"`
}

type Storage interface {
	CreateUser(ctx context.Context, username string, password string) (int, error)
	GetUserByUsername(ctx context.Context, username string) (int, error)

	Ping(ctx context.Context) error
	Close()
}

func GetStorage(ctx context.Context) (Storage, error) {
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage(ctx)
	} else {
		return nil, errors.New("empty db connection string")
	}
}
