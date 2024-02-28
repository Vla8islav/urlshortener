package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
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
	GetAllUserURLS(ctx context.Context, userID int) ([]storage.URLPair, error)
	GetShortenedURL(ctx context.Context, urlToShorten string) (string, int, error)
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

func (u URLShortenService) GetShortenedURL(ctx context.Context, urlToShorten string) (string, int, error) {
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
			return "", -1, fmt.Errorf("Couldn't generate shortened URL" + err.Error())
		}
		u.Storage.AddURLPair(ctx, newShortenedURL, urlToShorten, uuid.New().String(), auth.DefaultUserID)
		shortenedURL = newShortenedURL
	}
	return shortenedURL, auth.DefaultUserID, err
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

func (u URLShortenService) GetAllUserURLS(ctx context.Context, userID int) ([]storage.URLPair, error) {
	records, err := u.Storage.GetAllURLRecordsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return records, nil
}
