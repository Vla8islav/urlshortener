package storage

import "sync"

var mu sync.Mutex

func GetMakeshiftStorageInstance() *MakeshiftStorage {
	instance := new(MakeshiftStorage)
	instance.urlToShort = make(map[string]string)
	instance.shortToURL = make(map[string]string)

	return instance
}

type Storage interface {
	AddURLPair(shortenedURL string, fullURL string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
}

type MakeshiftStorage struct {
	urlToShort map[string]string
	shortToURL map[string]string
}

func (s MakeshiftStorage) AddURLPair(shortenedURL string, fullURL string) {
	mu.Lock()
	defer mu.Unlock()
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
}

func (s MakeshiftStorage) GetFullURL(shortenedURL string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.shortToURL[shortenedURL]
	return value, exists
}

func (s MakeshiftStorage) GetShortenedURL(fullURL string) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := s.urlToShort[fullURL]
	return value, exists
}
