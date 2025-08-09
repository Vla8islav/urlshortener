package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/helpers"
	"io"
	"net/http"
)

func (h *Handler) ShortenHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests are allowed to this endpoint", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "couldn't read request body", http.StatusBadRequest)
		return
	}

	var reqPayload ShortenRequestPayload
	err = json.Unmarshal(body, &reqPayload)
	if err != nil {
		http.Error(res, "couldn't parse the json request "+err.Error(), http.StatusBadRequest)
		return
	}

	fullURL, err := h.Service.GetByFull(reqPayload.FullURL)
	if err != nil {
		http.Error(res, "couldn't shorten url "+err.Error(), http.StatusInternalServerError)
		return
	}
	fullShortenedURL, err := helpers.GetFullShortenedURL(fullURL.ShortenedURL)
	if err != nil {
		http.Error(res, "couldn't generate full shortened url "+err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(ShortenResponsePayload{ShortenedURL: fullShortenedURL})
	if err != nil {
		http.Error(res, "couldn't shorten url "+err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-Type", "application/json")
	res.Header().Add("Content-Length", fmt.Sprintf("%d", len(resp)))
	_, err = res.Write(resp)
	if err != nil {
		http.Error(res, "couldn't write a response "+err.Error(), http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}
