package handlers

import (
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"net/http"
	"strings"
)

func GetUserURLSHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			http.Error(res, "Only GET requests are allowed to /api/user/urls", http.StatusBadRequest)
			return
		}

		authBearerStr := req.Header.Get("Authorization")
		if authBearerStr == "" {
			http.Error(res, "Needs Authorization header with JWT bearer to function",
				http.StatusBadRequest)
			return
		}
		bearerStr := strings.Replace(authBearerStr, "Bearer ", "", 1)
		userID, err := auth.GetUserID(bearerStr)
		if err != nil {
			http.Error(res, "Couldn't get user id from bearer",
				http.StatusBadRequest)
			return
		}

		urls, err := short.GetAllUserURLS(req.Context(), userID)

		responseBuffer, err := json.Marshal(urls)

		if err != nil {
			http.Error(res, "Failed to pack short url into json",
				http.StatusInternalServerError)
			return
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(responseBuffer)
	}
}
