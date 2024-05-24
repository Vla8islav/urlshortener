package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/api"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/errcustom"
	"io"
	"net/http"
)

func RegisterHandler(short api.GoBooruService) http.HandlerFunc {
	if short == (api.GoBooruService{}) {
		panic("Service wasn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			http.Error(res, "Only POST requests are allowed to /", http.StatusBadRequest)
			return
		}

		if req.Header.Get("Content-Type") != "application/json" {
			http.Error(res, "Content type must be application/json", http.StatusBadRequest)
			return
		}
		type CreateUserRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var requestStruct CreateUserRequest

		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Failed to read the request body", http.StatusInternalServerError)
			return

		}

		err = json.Unmarshal(body, &requestStruct)
		if err != nil {
			http.Error(res, "Incorrect json, expected json with username and password "+err.Error(), http.StatusBadRequest)
			return
		}

		newUserID, err := short.Register(req.Context(), requestStruct.Username, requestStruct.Password)
		if errors.Is(err, errcustom.ErrUserAlreadyExists) {
			http.Error(res, "user with the nickname "+requestStruct.Username+" already exist", http.StatusConflict)
			return
		} else if err != nil {
			http.Error(res, "couldn't create a user "+err.Error(), http.StatusInternalServerError)
			return
		}

		cookieValue, err := auth.BuildJWTString(newUserID)
		if err != nil {
			panic(fmt.Sprintf("Couldn't build a token string out of user %d %s", newUserID, err.Error()))
		}

		cookieName := "userid"
		cookie := http.Cookie{Name: cookieName,
			Value: cookieValue,
		}

		http.SetCookie(res, &cookie)
		res.WriteHeader(http.StatusCreated)
	}
}
