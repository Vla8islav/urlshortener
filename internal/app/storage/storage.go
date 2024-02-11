package storage

type Storage interface {
	AddURLPair(shortenedURL string, fullURL string, uuidStr string)
	AddURLPairInMemory(shortenedURL string, fullURL string, uuidStr string)
	GetFullURL(shortenedURL string) (string, bool)
	GetShortenedURL(fullURL string) (string, bool)
}
