package app

import (
	"context"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/url"
)

type GoBooruServiceMethods interface {
	SaveImage(ctx context.Context) ([]storage.URLPair, error)
}

type GoBooruService struct {
	Storage storage.Storage
}

func NewGoBooruService(ctx context.Context, s storage.Storage) (GoBooruService, error) {
	return GoBooruService{Storage: s}, nil
}

func (u GoBooruService) GetFullURL(ctx context.Context, shortenedPostfix string) (string, error) {
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := u.Storage.GetFullURL(ctx, fullSortURL)
	return longURL, found
}
