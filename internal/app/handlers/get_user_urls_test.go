package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetUserURLSHandler(t *testing.T) {
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
			name: "Successful link generation",
			request: func() *http.Request {
				u, err := url.Parse(randomShortenedURL)
				if err != nil {
					panic(err)
				}

				validRequest := httptest.NewRequest(http.MethodGet, u.Path, nil)
				validRequest.Header = http.Header{
					"Content-Type":  []string{"text/plain"},
					"Authorization": []string{"Bearer " + randomLinkBearer},
				}
				return validRequest

			},
			want: expectedResult{code: 200},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			GetUserURLSHandler(short)(w, testData.request())

			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)

		})
	}

}
