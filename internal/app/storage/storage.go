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
	Close()
}

func GetStorage() (Storage, error) {
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage(context.Background())
	}
	return NewMakeshiftStorage()
}
