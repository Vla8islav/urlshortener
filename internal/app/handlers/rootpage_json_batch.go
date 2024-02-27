package handlers

import (
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"io"
	"net/http"
)

func RootPageJSONBatchHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
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
		type URLShortenRequestRecord struct {
			CorrelationID string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		var requestStruct []URLShortenRequestRecord

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

		type URLShortenResponse struct {
			CorrelationID string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}
		var responseStruct []URLShortenResponse
		for _, record := range requestStruct {
			if !helpers.CheckIfItsURL(record.OriginalURL) {
				http.Error(res, "Incorrect url format", http.StatusBadRequest)
				return
			}

			shortenedURL, _ := short.GetShortenedURL(req.Context(), record.OriginalURL)

			responseStruct = append(responseStruct, URLShortenResponse{
				CorrelationID: record.CorrelationID,
				ShortURL:      shortenedURL,
			})

		}

		responseBuffer, err := json.Marshal(responseStruct)

		if err != nil {
			http.Error(res, "Failed to pack short url into json",
				http.StatusInternalServerError)
			return
		}

		res.Header().Add("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(responseBuffer)
	}

}
