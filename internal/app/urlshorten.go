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

type URLShortenService struct {
	s storage.Storage
}

type URLShortenServiceMethods interface {
	GetShortenedURL(urlToShorten string) string
	GetFullURL(shortenedPostfix string) (string, error)
	GenerateShortenedURL() (string, error)
	MatchesGeneratedURLFormat(s string) bool
}

func NewURLShortenService() URLShortenServiceMethods {
	return URLShortenService{s: storage.NewMakeshiftStorage()}
}

func (u URLShortenService) GetShortenedURL(urlToShorten string) string {
	if u.s == nil {
		panic("Database not initialised")
	}
	//s := storage.NewMakeshiftStorage()
	shortenedURL := ""
	if existingShortenedURL, alreadyExist := u.s.GetShortenedURL(urlToShorten); alreadyExist {
		shortenedURL = existingShortenedURL
	} else {
		newShortenedURL, err := u.GenerateShortenedURL()
		if err != nil {
			return ""
		}
		u.s.AddURLPair(newShortenedURL, urlToShorten)
		shortenedURL = newShortenedURL
	}
	return shortenedURL
}

var ErrURLNotFound = errors.New("couldn't find a requested URL")

func (u URLShortenService) GetFullURL(shortenedPostfix string) (string, error) {
	//s := storage.NewMakeshiftStorage()
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := u.s.GetFullURL(fullSortURL)
	if found {
		return longURL, nil
	} else {
		return longURL, ErrURLNotFound
	}
}

func (u URLShortenService) GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL,
		helpers.GenerateString(len(GeneratedShortenedURLSample), AllowedSymbolsInShortnedURL))
	if err != nil {
		return fullPath, err
	}
	return fullPath, nil
}

func (u URLShortenService) MatchesGeneratedURLFormat(s string) bool {
	s = strings.Trim(s, "/")
	r, _ := regexp.Compile("^[" + AllowedSymbolsInShortnedURL + "]+$")
	return len(s) == len(GeneratedShortenedURLSample) && r.MatchString(s)
}
