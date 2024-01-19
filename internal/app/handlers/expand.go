package handlers

import (
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app"
	"net/http"
)

type Handlers interface {
	ExpandHandler(res http.ResponseWriter, req *http.Request)
}

func ExpandHandler(short *app.URLShortenService) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			http.Error(res, "Only GET requests are allowed to /{id}", http.StatusBadRequest)
			return
		}

		uri := req.RequestURI
		if short.MatchesGeneratedURLFormat(uri) {
			fullURL, err := short.GetFullURL(uri)

			if err == nil {
				res.Header().Add("Location", fullURL)
				res.WriteHeader(http.StatusTemporaryRedirect)
			} else if errors.Is(err, app.ErrURLNotFound) {
				http.Error(res, "URL not found", http.StatusNotFound)
			} else {
				http.Error(res, "problem occured while extracting URL: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(res, "Invalid url format", http.StatusBadRequest)
		}
	}
}
