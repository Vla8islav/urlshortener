package handlers

import (
	"errors"
	helpers2 "github.com/Vla8islav/urlshortener/internal/helpers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/errlist"
	"net/http"
	"strings"
)

func (h *Handler) ExpandHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Only GET requests are allowed to /{id}", http.StatusBadRequest)
		return
	}

	uri := strings.Trim(req.RequestURI, "/")
	if helpers2.MatchesGeneratedURLFormat(uri) {
		urlObj, err := h.Service.GetByShortened(uri)
		switch {
		case err == nil:
			res.Header().Add("Location", urlObj.FullURL)
			res.WriteHeader(http.StatusTemporaryRedirect)
		case errors.Is(err, errlist.ErrURLNotFound):
			http.Error(res, "URLRepo not found", http.StatusNotFound)
		default:
			http.Error(res, "problem occured while extracting URLRepo: "+err.Error(),
				http.StatusInternalServerError)
		}

	} else {
		http.Error(res, "Invalid url format", http.StatusBadRequest)
	}
}
