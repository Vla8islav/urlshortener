package db

import (
	"github.com/Vla8islav/urlshortener/internal/domain"
	"github.com/Vla8islav/urlshortener/internal/helpers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/errlist"
	"sync"
)

type MakeshiftInMemoryDB struct {
	fullURLToShort map[string]string
	shortToFullURL map[string]string
}

var instance *MakeshiftInMemoryDB = nil

func GetInstance() *MakeshiftInMemoryDB {
	sync.OnceFunc(func() {
		instance = new(MakeshiftInMemoryDB)
		instance.fullURLToShort = make(map[string]string)
		instance.shortToFullURL = make(map[string]string)
	})()

	return instance
}

func (s MakeshiftInMemoryDB) GetByFull(fullURL string) (domain.URL, error) {
	shortURL, exists := s.fullURLToShort[fullURL]
	var err error
	if !exists {
		shortURL = helpers.GenerateShortenedURLUID()
		s.fullURLToShort[fullURL] = shortURL
		s.shortToFullURL[shortURL] = fullURL

	}
	return domain.URL{
		FullURL:      fullURL,
		ShortenedURL: shortURL,
	}, err
}

func (s MakeshiftInMemoryDB) GetByShortened(shortURL string) (domain.URL, error) {
	fullURL, exists := s.shortToFullURL[shortURL]
	if !exists {
		return domain.URL{}, errlist.ErrURLNotFound

	}
	return domain.URL{
		FullURL:      fullURL,
		ShortenedURL: shortURL,
	}, nil
}
