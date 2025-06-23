package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/application"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootPageHandler(t *testing.T) {
	repo := db.GetInstance()
	service := application.NewURLShortenService(repo)
	handler := NewHandler(service)

	type expectedResult struct {
		code        int
		contentType string
	}

	validRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	validRequest.Header = http.Header{
		"Content-Type": []string{"text/plain; charset=utf-8"},
	}
	validRequest.Body = io.NopCloser(strings.NewReader("http://ya.ru"))

	getRequest := httptest.NewRequest(http.MethodGet, "/", nil)

	testDataArray := []struct {
		name    string
		request *http.Request
		want    expectedResult
	}{
		{
			name:    "Successful link generation",
			request: validRequest,
			want: expectedResult{
				code: 201,
			},
		},
		{
			name:    "Successful link generation",
			request: getRequest,
			want:    expectedResult{code: 400},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handler.RootPageHandler(w, testData.request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			regexToValidateTheLink := strings.TrimRight(config.ReadFlags().ShortenerBaseURL, "/") + "/[a-zA-Z]{8}"
			if w.Code >= 200 && w.Code <= 299 {
				assert.Regexp(t, regexToValidateTheLink, string(resBody))
				assert.Equal(t, testData.want.contentType, res.Header.Get("Content-Type"))
			}

		})
	}

}
