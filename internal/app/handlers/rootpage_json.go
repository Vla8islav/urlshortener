package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"io"
	"net/http"
)

func RootPageJSONHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Underlying infrastructure isn't initialised")
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
		type URLShortenRequest struct {
			URL string `json:"url"`
		}

		var requestStruct URLShortenRequest

		body, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, "Failed to read the request body", http.StatusInternalServerError)
			return

		}

		err = json.Unmarshal(body, &requestStruct)
		if err != nil {
			http.Error(res, "Incorrect json", http.StatusBadRequest)
			return
		}

		if !helpers.CheckIfItsURL(requestStruct.URL) {
			http.Error(res, "Incorrect url format", http.StatusBadRequest)
			return
		}

		shortenedURL, shortURLError := short.GetShortenedURL(requestStruct.URL)

		returnStatus := http.StatusCreated
		var urlAlreadyExist *app.UrlExistError
		if errors.As(shortURLError, &urlAlreadyExist) {
			returnStatus = http.StatusConflict
		}

		type URLShortenResponse struct {
			Result string `json:"result"`
		}

		responseStruct := URLShortenResponse{Result: shortenedURL}
		responseBuffer, err := json.Marshal(responseStruct)

		if err != nil {
			http.Error(res, "Failed to pack short url '"+shortenedURL+"' into json",
				http.StatusInternalServerError)
			return
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(returnStatus)
		res.Write(responseBuffer)
	}

}
