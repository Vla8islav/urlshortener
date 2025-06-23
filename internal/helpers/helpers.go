package helpers

import "net/url"

func CheckIfItsURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}
