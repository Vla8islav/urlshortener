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
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
	mu.Unlock()
}

func (s MakeshiftStorage) GetFullURL(shortenedURL string) (string, bool) {
	mu.Lock()
	value, exists := s.shortToURL[shortenedURL]
	mu.Unlock()
	return value, exists
}

func (s MakeshiftStorage) GetShortenedURL(fullURL string) (string, bool) {
	mu.Lock()
	value, exists := s.urlToShort[fullURL]
	mu.Unlock()
	return value, exists
}
