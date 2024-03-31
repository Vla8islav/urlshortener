package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/concurrency"
	"io"
	"net/http"
)

func UserURLSHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		var userID int
		var err error

		bearerHeader := req.Header.Get("Authorization")
		if bearerHeader != "" {
			bearer := auth.GetBearerFromBearerHeader(bearerHeader)
			userID, err = auth.GetUserID(bearer)
			if err != nil {
				http.Error(res, "Needs Authorization header with JWT bearer to function",
					http.StatusUnauthorized)
				return

			}
		} else {
			cookieName := "userid"
			existingCookie, err := req.Cookie(cookieName)
			if errors.Is(err, http.ErrNoCookie) {
				http.Error(res, "Needs Authorization cookie with JWT bearer to function",
					http.StatusUnauthorized)
				return
			}

			userID, err = auth.GetUserID(existingCookie.Value)
			if err != nil {
				http.Error(res, "Couldn't get user id from bearer",
					http.StatusBadRequest)
				return
			}
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

			queue := concurrency.NewQueue()

			for i := 0; i < len(urls); i++ {
				w := concurrency.NewWorker(i, queue, concurrency.NewDeleter(&short, context.Background()))
				go w.Loop()
			}

			for _, url := range urls {
				queue.Push(&concurrency.Task{URL: url, UserID: userID})
			}
			res.WriteHeader(http.StatusAccepted)

		default:
			http.Error(res, "Only GET or DELETE requests are allowed to /api/user/urls", http.StatusBadRequest)
			return
		}

	}
}
