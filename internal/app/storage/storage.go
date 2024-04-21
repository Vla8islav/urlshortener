package storage

import (
	"context"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
)

type URLPair struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
	Deleted  bool   `json:"deleted"`
}

type Storage interface {
	AddURLPair(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int)
	AddURLPairInMemory(ctx context.Context, shortenedURL string, fullURL string, uuidStr string, userID int)
	GetFullURL(ctx context.Context, shortenedURL string) (string, error)
	GetShortenedURL(ctx context.Context, fullURL string) (string, int, bool)

	GetAllURLRecordsByUser(ctx context.Context, userID int) ([]URLPair, error)
	GetNewUserID(ctx context.Context) (int, error)

	DeleteURL(ctx context.Context, shortenedURL string, userID int) error

	Ping(ctx context.Context) error
	Close()
}

func GetStorage(ctx context.Context) (Storage, error) {
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage(ctx)
	}
	return NewMakeshiftStorage(ctx)
}
