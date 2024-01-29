package handlers

import (
	"encoding/json"
	"fmt"
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

		shortenedURL := short.GetShortenedURL(requestStruct.URL)

		type URLShortenResponse struct {
			Result string `json:"result"`
		}

		responseStruct := URLShortenResponse{Result: shortenedURL}
		responseBuffer, err := json.Marshal(responseStruct)
		// TODO: see why this is necessary for compression to function
		responseBufferStr := string(responseBuffer)

		if err != nil {
			http.Error(res, "Failed to pack short url '"+shortenedURL+"' into json",
				http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Header().Add("Content-Type", "application/json")
		res.Header().Add("Content-Length", fmt.Sprintf("%d", len(responseBuffer)))
		res.Write([]byte(responseBufferStr))
	}

}
