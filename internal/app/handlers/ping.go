package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/http"
)

func PingHandler(s *storage.Storage) http.HandlerFunc {
	if s == nil {
		panic("Underlying storage isn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {
		err := (*s).Ping(req.Context())
		if err == nil {
			res.WriteHeader(http.StatusOK)
		} else {
			res.Write([]byte(err.Error()))
			res.WriteHeader(http.StatusInternalServerError)
		}
	}

}
