package storage

import "github.com/Vla8islav/urlshortener/internal/app/configuration"

type Storage interface {
	AddURLPair(shortenedURL string, fullURL string, uuidStr string)
	AddURLPairInMemory(shortenedURL string, fullURL string, uuidStr string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
}

func GetStorage() (Storage, error) {
	//s, err := storage.NewMakeshiftStorage()
	if configuration.ReadFlags().DBConnectionString != "" {
		return NewPostgresStorage()
	}
	return NewMakeshiftStorage()
}
