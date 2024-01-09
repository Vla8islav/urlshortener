package storage

var urlToShort = map[string]string{}
var shortToURL = map[string]string{}

func AddURLPair(shortenedURL string, fullURL string) {
	urlToShort[fullURL] = shortenedURL
	shortToURL[shortenedURL] = fullURL
}

func GetFullURL(shortenedURL string) (string, bool) {
	value, exists := shortToURL[shortenedURL]
	return value, exists
}

func GetShortenedURL(fullURL string) (string, bool) {
	value, exists := urlToShort[fullURL]
	return value, exists
}
