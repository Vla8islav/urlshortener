package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestExpandHandler(t *testing.T) {
	s, err := storage.GetStorage()
	if err != nil {
		panic(err)
	}
	short, _ := app.NewURLShortenService(s)
	shortenedURL, _ := short.GetShortenedURL("http://ya.ru")

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
				u, err := url.Parse(shortenedURL)
				if err != nil {
					panic(err)
				}

				validRequest := httptest.NewRequest(http.MethodGet, u.Path, nil)
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

			ExpandHandler(short)(w, testData.request())

			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)
			//resBody, err := io.ReadAll(res.Body)
			//
			//require.NoError(t, err)
			//if w.Code >= 200 && w.Code <= 299 {
			//	assert.JSONEq(t, testData.want.response, string(resBody))
			//	assert.Equal(t, testData.want.contentType, res.Header.Get("Content-Type"))
			//}

		})
	}

}
