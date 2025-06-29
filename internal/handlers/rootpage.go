package handlers

import (
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/helpers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"io"
	"net/http"
)

func (h *Handler) RootPageHandler(res http.ResponseWriter, req *http.Request) {
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

	shortenedURLObj, err := h.Service.GetByFull(bodyString)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		http.Error(res, "problem occured while trying to fetch shortened url: "+err.Error(), http.StatusInternalServerError)

		return
	}

	shortenedURL := config.ReadFlags().ShortenerBaseURL + "/" + shortenedURLObj.ShortenedURL
	res.WriteHeader(http.StatusCreated)
	res.Header().Add("Content-Type", "text/plain")
	res.Header().Add("Content-Length", fmt.Sprintf("%d", len(shortenedURL)))
	_, err = res.Write([]byte(shortenedURL))
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		http.Error(res, "problem occured while writing url: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
