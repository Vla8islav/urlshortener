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

func DeleteUserURLSHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodDelete {
			http.Error(res, "Delete url user routing failue. Expected DELETE, got "+req.Method,
				http.StatusInternalServerError)
			return

		}

		var userID int
		var err error

		cookieName := "userid"
		existingCookie, err := req.Cookie(cookieName)
		if errors.Is(err, http.ErrNoCookie) {
			http.Error(res, "needs Authorization cookie with JWT bearer to function "+err.Error(),
				http.StatusUnauthorized)
			return
		}

		userID, err = auth.GetUserID(existingCookie.Value)
		if err != nil {
			http.Error(res, "Couldn't get user id from bearer",
				http.StatusBadRequest)
			return
		}

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

	}
}
