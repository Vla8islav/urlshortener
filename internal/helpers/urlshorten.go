package helpers

import (
	"errors"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"net/url"
	"regexp"
	"strings"
)

const AllowedSymbolsInShortnedURL = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const GeneratedShortenedURLSample = "EwHXdJfB"
const ShortenedURLLength = len(GeneratedShortenedURLSample)

func GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath(config.ReadFlags().ShortenerBaseURL, GenerateShortenedURLUID())
	if err != nil {
		return fullPath, err
	}
	return fullPath, nil
}

func GenerateShortenedURLUID() string {
	return GenerateString(ShortenedURLLength, AllowedSymbolsInShortnedURL)
}

func MatchesGeneratedURLFormat(s string) bool {
	s = strings.Trim(s, "/")
	r, _ := regexp.Compile("^[" + AllowedSymbolsInShortnedURL + "]+$")
	return len(s) == ShortenedURLLength && r.MatchString(s)
}

func GetFullShortenedURL(shortenedURL string) (string, error) {
	if !MatchesGeneratedURLFormat(shortenedURL) {
		return "", errors.New("shortened URL does not match the expected format")
	}
	return url.JoinPath(config.ReadFlags().ShortenerBaseURL, shortenedURL)
}

func GetFullShortenedURLSample() string {
	sample, err := GetFullShortenedURL(GeneratedShortenedURLSample)
	if err != nil {
		panic("Failed to generate full shortened URL sample: " + err.Error())
	}
	return sample
}
