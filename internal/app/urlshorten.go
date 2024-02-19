package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/google/uuid"
	"net/url"
	"regexp"
	"strings"
)

const AllowedSymbolsInShortnedURL = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const GeneratedShortenedURLSample = "EwHXdJfB"

type URLShortenService struct {
	Storage storage.Storage
}

type URLShortenServiceMethods interface {
	GetShortenedURL(ctx context.Context, urlToShorten string) (string, error)
	GetFullURL(ctx context.Context, shortenedPostfix string) (string, error)
	GenerateShortenedURL(ctx context.Context) (string, error)
	MatchesGeneratedURLFormat(s string) bool
}

func NewURLShortenService(ctx context.Context, s storage.Storage) (URLShortenServiceMethods, error) {

	return URLShortenService{Storage: s}, nil
}

type URLExistError struct {
	URL string
	Err error
}

func (ue *URLExistError) Error() string {
	return fmt.Sprintf("URL: %s Error: %v", ue.URL, ue.Err)
}

func (u URLShortenService) GetShortenedURL(ctx context.Context, urlToShorten string) (string, error) {
	if u.Storage == nil {
		panic("Database not initialised")
	}
	shortenedURL := ""
	var err error
	if existingShortenedURL, alreadyExist := u.Storage.GetShortenedURL(ctx, urlToShorten); alreadyExist {
		shortenedURL = existingShortenedURL
		err = &URLExistError{Err: err, URL: existingShortenedURL}
	} else {
		newShortenedURL, err := u.GenerateShortenedURL(ctx)
		if err != nil {
			return "", fmt.Errorf("Couldn't generate shortened URL" + err.Error())
		}
		u.Storage.AddURLPair(ctx, newShortenedURL, urlToShorten, uuid.New().String())
		shortenedURL = newShortenedURL
	}
	return shortenedURL, err
}

var ErrURLNotFound = errors.New("couldn't find a requested URL")

func (u URLShortenService) GetFullURL(ctx context.Context, shortenedPostfix string) (string, error) {
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := u.Storage.GetFullURL(ctx, fullSortURL)
	if found {
		return longURL, nil
	} else {
		return longURL, ErrURLNotFound
	}
}

func (u URLShortenService) GenerateShortenedURL(ctx context.Context) (string, error) {
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
