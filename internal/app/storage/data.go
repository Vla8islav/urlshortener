package storage

import "sync"

var instance *makeshiftStorage = nil

func GetInstance() MakeshiftStorage {
	sync.OnceFunc(func() {
		if instance == nil {
			instance = new(makeshiftStorage)
			instance.urlToShort = make(map[string]string)
			instance.shortToURL = make(map[string]string)
		}
	})()

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
	s.urlToShort[fullURL] = shortenedURL
	s.shortToURL[shortenedURL] = fullURL
}

func (s makeshiftStorage) GetFullURL(shortenedURL string) (string, bool) {
	value, exists := s.shortToURL[shortenedURL]
	return value, exists
}

func (s makeshiftStorage) GetShortenedURL(fullURL string) (string, bool) {
	value, exists := s.urlToShort[fullURL]
	return value, exists
}
