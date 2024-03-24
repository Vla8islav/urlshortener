package helpers

import (
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"net/url"
	"strings"
)

func CheckIfItsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

func ShortKeyToURL(shortKey string) (string, error) {
	return url.JoinPath(configuration.ReadFlags().ShortenerBaseURL, shortKey)
}

func URLToShortKey(URL string) string {
	return strings.TrimLeft(URL, configuration.ReadFlags().ShortenerBaseURL)
}
