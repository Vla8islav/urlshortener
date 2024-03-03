package handlers

import (
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"io"
	"net/http"
)

func UserURLSHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		switch req.Method {
		case http.MethodGet:

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
			buffer, err := io.ReadAll(req.Body)
			if err != nil {
				http.Error(res, "Couldn't read DELETE request body "+err.Error(),
					http.StatusInternalServerError)
				return
			}
			var urls []string
			err = json.Unmarshal(buffer, &urls)
			if err != nil {
				http.Error(res, "Incorrect format of the DELETE request body "+err.Error(),
					http.StatusInternalServerError)
				return
			}

			for _, url := range urls {
				err = short.DeleteLink(req.Context(), url)
				if err != nil {
					http.Error(res, "Couldn't delete link/userID pair "+url+err.Error(),
						http.StatusInternalServerError)
					return
				}

			}
			res.WriteHeader(http.StatusAccepted)

		default:
			http.Error(res, "Only GET or DELETE requests are allowed to /api/user/urls", http.StatusBadRequest)
			return
		}

	}
}
