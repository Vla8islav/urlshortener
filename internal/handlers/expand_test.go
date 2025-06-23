package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/application"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExpandHandler(t *testing.T) {

	repo := db.GetInstance()
	service := application.NewURLShortenService(repo)
	handler := NewHandler(service)

	shortenedURL, _ := handler.Service.GetShortURL("http://ya.ru")

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

				validRequest := httptest.NewRequest(http.MethodGet, "/"+shortenedURL.ShortenedURL, nil)
				validRequest.Header = http.Header{
					"Content-Type": []string{"text/plain"},
				}
				return validRequest

			},
			want: expectedResult{code: 307},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handler.ExpandHandler(w, testData.request())

			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)

		})
	}

}
