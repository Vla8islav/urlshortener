package app

import (
	"context"
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
	GetShortenedURL(ctx context.Context, urlToShorten string, bearerToken string) (string, int, error)
	GetFullURL(ctx context.Context, shortenedPostfix string) (string, error)
	GenerateShortenedURL(ctx context.Context) (string, error)
	MatchesGeneratedURLFormat(s string) bool
	DeleteLink(ctx context.Context, shortenedURL string, userID int) error
	//DeleteLink(ctx context.Context, shortenedURL string, userID int) error
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

func (u URLShortenService) GetShortenedURL(ctx context.Context,
	urlToShorten string,
	bearerHeader string) (string, int, error) {
	if u.Storage == nil {
		panic("Database not initialised")
	}

	bearer := auth.GetBearerFromBearerHeader(bearerHeader)

	var err error
	userID, err := auth.GetUserID(bearer)
	if err != nil {
		userID = -1
	}

	shortenedURL := ""

	if existingShortenedURL, id, alreadyExist := u.Storage.GetShortenedURL(ctx, urlToShorten); alreadyExist {
		shortenedURL = existingShortenedURL
		err = &URLExistError{Err: err, URL: existingShortenedURL}
		userID = id
	} else {
		if userID == -1 {
			userID, err = u.Storage.GetNewUserID(ctx)
		}
		if err != nil {
			return "", -1, fmt.Errorf("couldn't create new user id" + err.Error())
		}
		newShortenedURL, err := u.GenerateShortenedURL(ctx)
		if err != nil {
			return "", -1, fmt.Errorf("Couldn't generate shortened URL" + err.Error())
		}
		u.Storage.AddURLPair(ctx, newShortenedURL, urlToShorten, uuid.New().String(), userID)
		shortenedURL = newShortenedURL
	}
	return shortenedURL, userID, err
}

func (u URLShortenService) GetFullURL(ctx context.Context, shortenedPostfix string) (string, error) {
	fullSortURL, err := url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortenedPostfix)
	if err != nil {
		return "", err
	}
	longURL, found := u.Storage.GetFullURL(ctx, fullSortURL)
	return longURL, found
}

func (u URLShortenService) GenerateShortenedURL(ctx context.Context) (string, error) {
	shortKey := helpers.GenerateString(len(GeneratedShortenedURLSample), AllowedSymbolsInShortnedURL)
	fullPath, err := helpers.ShortKeyToURL(shortKey)
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

func (u URLShortenService) DeleteLink(ctx context.Context, shortenedURL string, userID int) error {
	return u.Storage.DeleteURL(ctx, shortenedURL, userID)
}
