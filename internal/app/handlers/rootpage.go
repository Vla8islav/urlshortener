package handlers

import (
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"io"
	"net/http"
)

func RootPageHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed to /", http.StatusBadRequest)
		return
	}

	if req.Header.Get("Content-Type") != "text/plain; charset=utf-8" {
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

	shortenedURL := app.GetShortenedURL(bodyString)

	res.WriteHeader(http.StatusCreated)
	res.Header().Add("Content-Type", "text/plain")
	res.Header().Add("Content-Length", fmt.Sprintf("%d", len(shortenedURL)))
	res.Write([]byte(shortenedURL))
}
