package storage

import (
	"context"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
)

type Storage interface {
	AddURLPair(shortenedURL string, fullURL string, uuidStr string)
	AddURLPairInMemory(shortenedURL string, fullURL string, uuidStr string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
	Ping() error
	Close()
}

func GetStorage(ctx context.Context) (Storage, error) {
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage(ctx)
	}
	return NewMakeshiftStorage()
}
