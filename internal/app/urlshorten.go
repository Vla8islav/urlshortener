package app

import (
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/url"
	"regexp"
	"strings"
)

const AllowedSymbolsInShortnedURL = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const GeneratedShortenedURLSample = "EwHXdJfB"

func GetShortenedURL(urlToShorten string) string {
	s := storage.GetInstance()
	shortenedURL := ""
	if existingShortenedURL, alreadyExist := s.GetShortenedURL(urlToShorten); alreadyExist {
		shortenedURL = existingShortenedURL
	} else {
		newShortenedURL, err := GenerateShortenedURL()
		if err != nil {
			return ""
		}
		s.AddURLPair(newShortenedURL, urlToShorten)
		shortenedURL = newShortenedURL
	}
	return shortenedURL
}

var ErrURLNotFound = errors.New("couldn't find a requested URL")

func GetFullURL(shortenedPostfix string) (string, error) {
	s := storage.GetInstance()
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := s.GetFullURL(fullSortURL)
	if found {
		return longURL, nil
	} else {
		return longURL, ErrURLNotFound
	}
}

func GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL,
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
