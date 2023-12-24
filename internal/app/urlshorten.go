package app

import (
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/url"
	"regexp"
	"strings"
)

const AllowedSymbolsInShortnedURL = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const GeneratedShortenedURLSample = "EwHXdJfB"

func GetShortenedURL(urlToShorten string) string {
	shortenedURL, err := GenerateShortenedURL()
	if err != nil {
		return ""
	}
	storage.AddURLPair(shortenedURL, urlToShorten)
	return shortenedURL
}

func GetFullURL(shortenedPostfix string) string {
	fullSortURL, err := url.JoinPath("http://localhost:8080/", shortenedPostfix)
	if err != nil {
		return ""
	}
	longURL, _ := storage.GetFullURL(fullSortURL)
	return longURL
}

func GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath("http://localhost:8080/",
		helpers.GenerateString(len(GeneratedShortenedURLSample), AllowedSymbolsInShortnedURL))
	if err != nil {
		return fullPath, err
	}
	return fullPath, nil
}

func MatchesGeneratedURLFormat(s string) bool {
	s = strings.Trim(s, "/")
	r, _ := regexp.Compile("^[" + AllowedSymbolsInShortnedURL + "]+$")
	return len(s) == len(GeneratedShortenedURLSample) && r.MatchString(s)
}
