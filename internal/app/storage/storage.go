package storage

import (
	"context"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
)

type Storage interface {
	AddURLPair(ctx context.Context, shortenedURL string, fullURL string, uuidStr string)
	AddURLPairInMemory(ctx context.Context, shortenedURL string, fullURL string, uuidStr string)
	GetFullURL(ctx context.Context, shortenedURL string) (string, bool)
	GetShortenedURL(ctx context.Context, fullURL string) (string, bool)
	Ping(ctx context.Context) error
	Close()
}

func GetStorage(ctx context.Context) (Storage, error) {
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage(ctx)
	}
	return NewMakeshiftStorage(ctx)
}
