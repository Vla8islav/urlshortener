package storage

import "sync"

var instance *makeshiftStorage = nil
var mu sync.Mutex

func GetInstance() MakeshiftStorage {
	mu.Lock()
	if instance == nil {
		instance = new(makeshiftStorage)
		instance.urlToShort = make(map[string]string)
		instance.shortToURL = make(map[string]string)
	}
	mu.Unlock()

	return instance
}

type MakeshiftStorage interface {
	AddURLPair(shortenedURL string, fullURL string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
}

type makeshiftStorage struct {
	urlToShort map[string]string
	shortToURL map[string]string
}

func (s makeshiftStorage) AddURLPair(shortenedURL string, fullURL string) {
	mu.Lock()
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
	mu.Unlock()
}

func (s makeshiftStorage) GetFullURL(shortenedURL string) (string, bool) {
	mu.Lock()
	value, exists := s.shortToURL[shortenedURL]
	mu.Unlock()
	return value, exists
}

func (s makeshiftStorage) GetShortenedURL(fullURL string) (string, bool) {
	mu.Lock()
	value, exists := s.urlToShort[fullURL]
	mu.Unlock()
	return value, exists
}
