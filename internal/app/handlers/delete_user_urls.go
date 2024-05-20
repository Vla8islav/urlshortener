package handlers

import (
	"context"
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"io"
	"net/http"
)

func DeleteUserURLSHandler(short app.GoBooruService) http.HandlerFunc {
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

		bearer := auth.GetBearerNewOrOld(res, req)

		userID, err = auth.GetUserID(bearer)
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

		for _, url := range urls {
			go short.DeleteLink(context.Background(), url, userID)
		}
		res.WriteHeader(http.StatusAccepted)

	}
}
