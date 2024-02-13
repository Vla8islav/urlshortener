package handlers

import (
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"io"
	"net/http"
	"strings"
)

func RootPageHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Underlying infrastructure isn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed to /", http.StatusBadRequest)
			return
		}

		if !(req.Header.Get("Content-Type") == "text/plain; charset=utf-8" ||
			strings.Contains(req.Header.Get("Content-Type"), "gzip")) {
			http.Error(res, "Content type must be text/plain", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Failed to read the request body", http.StatusInternalServerError)
			return

		}
		bodyString := string(body)
		if !helpers.CheckIfItsURL(bodyString) {
			http.Error(res, "Incorrect url format", http.StatusBadRequest)
			return
		}

		shortenedURL, shortURLError := short.GetShortenedURL(bodyString)

		returnStatus := http.StatusCreated
		var urlAlreadyExist *app.UrlExistError
		if errors.As(shortURLError, &urlAlreadyExist) {
			returnStatus = http.StatusConflict
		}

		res.WriteHeader(returnStatus)
		res.Header().Add("Content-Type", "text/plain")
		res.Write([]byte(shortenedURL))
	}

}
