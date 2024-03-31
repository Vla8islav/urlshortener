package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteUserURLSHandler(t *testing.T) {
	ctx, cancel := helpers.GetDefaultContext()
	defer cancel()

	s, err := storage.GetStorage(ctx)

	if err != nil {
		panic(err)
	}
	short, _ := app.NewURLShortenService(ctx, s)
	randomBaseURL := helpers.GenerateRandomURL()
	randomShortenedURL, userID, _ := short.GetShortenedURL(ctx, randomBaseURL, "")
	randomLinkBearer, _ := auth.BuildJWTString(userID)

	type expectedResult struct {
		code int
	}

	testDataArray := []struct {
		name    string
		request func() *http.Request
		want    expectedResult
	}{
		{
			name: "Successful deletion",
			request: func() *http.Request {
				body := []string{helpers.URLToShortKey(randomShortenedURL)}
				bodyJSON, err := json.Marshal(body)
				if err != nil {
					panic("Couldn't do a deletion test")
				}
				bodyReaderJSON := bytes.NewReader(bodyJSON)

				validRequest := httptest.NewRequest(http.MethodDelete, "/", bodyReaderJSON)
				validRequest.Header = http.Header{
					"Content-Type": []string{"text/plain"},
					"Cookie":       []string{"userid=" + randomLinkBearer},
				}
				return validRequest

			},
			want: expectedResult{code: http.StatusAccepted},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			DeleteUserURLSHandler(short)(w, testData.request())

			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)

		})
	}

}
