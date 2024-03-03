package handlers

import (
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"net/http"
)

func GetUserURLSHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		authBearerStr := req.Header.Get("Authorization")
		if authBearerStr == "" {
			http.Error(res, "Needs Authorization header with JWT bearer to function",
				http.StatusUnauthorized)
			return
		}
		bearerStr := auth.GetBearerFromBearerHeader(authBearerStr)
		userID, err := auth.GetUserID(bearerStr)
		if err != nil {
			http.Error(res, "Couldn't get user id from bearer",
				http.StatusBadRequest)
			return
		}

		switch req.Method {
		case http.MethodGet:
			urls, err := short.GetAllUserURLS(req.Context(), userID)
			if err != nil {
				http.Error(res, "Failed to get all user urls "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			responseBuffer, err := json.Marshal(urls)

			if err != nil {
				http.Error(res, "Failed to pack short url into json "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			res.Header().Add("Content-Type", "application/json")
			res.WriteHeader(http.StatusOK)
			res.Write(responseBuffer)

		case http.MethodDelete:

		default:
			http.Error(res, "Only GET or DELETE requests are allowed to /api/user/urls", http.StatusBadRequest)
			return
		}

	}
}
