package helpers

import (
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"net/url"
	"regexp"
	"strings"
)

const AllowedSymbolsInShortnedURL = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const GeneratedShortenedURLSample = "EwHXdJfB"

func GenerateShortenedURL() (string, error) {
	fullPath, err := url.JoinPath(config.ReadFlags().ShortenerBaseURL, GenerateShortenedUrlUid())
	if err != nil {
		return fullPath, err
	}
	return fullPath, nil
}

func GenerateShortenedUrlUid() string {
	return GenerateString(len(GeneratedShortenedURLSample), AllowedSymbolsInShortnedURL)
}

func MatchesGeneratedURLFormat(s string) bool {
	s = strings.Trim(s, "/")
	r, _ := regexp.Compile("^[" + AllowedSymbolsInShortnedURL + "]+$")
	return len(s) == len(GeneratedShortenedURLSample) && r.MatchString(s)
}
