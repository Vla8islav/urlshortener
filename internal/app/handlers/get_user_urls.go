package handlers

import (
	"context"
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/concurrency"
	"io"
	"net/http"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

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

			queue := concurrency.NewQueue()
			const MAX_WORKER_COUNT = 10

			for i := 0; i < min(MAX_WORKER_COUNT, len(urls)); i++ {
				w := concurrency.NewWorker(i, queue, concurrency.NewDeleter(&short, context.Background()))
				go w.Loop()
			}

			for _, url := range urls {
				queue.Push(&concurrency.Task{URL: url})
			}
			res.WriteHeader(http.StatusAccepted)

		default:
			http.Error(res, "Only GET or DELETE requests are allowed to /api/user/urls", http.StatusBadRequest)
			return
		}

	}
}
