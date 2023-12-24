package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"net/http"
)

func IDHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed to /{id}", http.StatusBadRequest)
		return
	}

	if req.Header.Get("Content-Type") != "text/plain" {
		http.Error(res, "Content type must be text/plain", http.StatusBadRequest)
		return
	}

	uri := req.RequestURI
	if app.MatchesGeneratedURLFormat(uri) {
		fullURL := app.GetFullURL(uri)
		if len(fullURL) > 0 {
			res.WriteHeader(http.StatusTemporaryRedirect)
			res.Write([]byte(fullURL))
		} else {
			http.Error(res, "Url not found", http.StatusNotFound)
		}
	} else {
		http.Error(res, "Invalid url format", http.StatusBadRequest)
	}
}
