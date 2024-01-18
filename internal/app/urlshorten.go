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

type URLShorten struct {
	S storage.Storage
}

func (u URLShorten) GetShortenedURL(urlToShorten string) string {
	if u.S == nil {
		panic("Database not initialised")
	}
	//S := storage.GetMakeshiftStorageInstance()
	shortenedURL := ""
	if existingShortenedURL, alreadyExist := u.S.GetShortenedURL(urlToShorten); alreadyExist {
		shortenedURL = existingShortenedURL
	} else {
		newShortenedURL, err := u.GenerateShortenedURL()
		if err != nil {
			return ""
		}
		u.S.AddURLPair(newShortenedURL, urlToShorten)
		shortenedURL = newShortenedURL
	}
	return shortenedURL
}

var ErrURLNotFound = errors.New("couldn't find a requested URL")

func (u URLShorten) GetFullURL(shortenedPostfix string) (string, error) {
	//S := storage.GetMakeshiftStorageInstance()
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := u.S.GetFullURL(fullSortURL)
	if found {
		return longURL, nil
	} else {
		return longURL, ErrURLNotFound
	}
}

func (u URLShorten) GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL,
		helpers.GenerateString(len(GeneratedShortenedURLSample), AllowedSymbolsInShortnedURL))
	if err != nil {
		return fullPath, err
	}
	return fullPath, nil
}

func (u URLShorten) MatchesGeneratedURLFormat(s string) bool {
	s = strings.Trim(s, "/")
	r, _ := regexp.Compile("^[" + AllowedSymbolsInShortnedURL + "]+$")
	return len(s) == len(GeneratedShortenedURLSample) && r.MatchString(s)
}
