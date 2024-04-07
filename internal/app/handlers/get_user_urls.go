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

		if req.Method != http.MethodGet {
			http.Error(res, "Get user urls routing failue. Expected GET, got "+req.Method,
				http.StatusInternalServerError)
			return

		}

		var userID int
		var err error

		bearer := auth.GetBearerNewOrOld(res, req)

		userID, err = auth.GetUserID(bearer)
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

	}
}
