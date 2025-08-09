package handlers

import (
	"encoding/json"
	"github.com/Vla8islav/urlshortener/internal/application"
	"github.com/Vla8islav/urlshortener/internal/helpers"
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

func TestShortenHandler(t *testing.T) {
	repo := db.GetInstance()
	service := application.NewURLShortenService(repo)
	handler := NewHandler(service)

	validRequest := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
	validRequest.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}
	shortenRequest, err := json.Marshal(ShortenRequestPayload{FullURL: "http://ya.ru"})
	if err != nil {
		panic("couldn't marshal request payload: " + err.Error())
	}
	validRequest.Body = io.NopCloser(strings.NewReader(string(shortenRequest)))
	getRequest := httptest.NewRequest(http.MethodGet, "/api/shorten", nil)

	type expectedResult struct {
		code        int
		contentType string
	}

	testDataArray := []struct {
		name    string
		request *http.Request
		want    expectedResult
	}{
		{
			name:    "Successful link generation",
			request: validRequest,
			want: expectedResult{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},
		{
			name:    "400 response for GET request",
			request: getRequest,
			want:    expectedResult{code: http.StatusBadRequest},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handler.ShortenHandler(w, testData.request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			if w.Code >= 200 && w.Code <= 299 {
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				// проверяем, что тело ответа соответствует ожидаемому
				assert.NotEmpty(t, resBody)
				// проверяем, что сгенерированная ссылка соответствует ожидаемому формату
				var resp ShortenResponsePayload
				err = json.Unmarshal(resBody, &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.ShortenedURL, len(helpers.GetFullShortenedURLSample()))

				regexToValidateTheLink := strings.TrimRight(config.ReadFlags().ShortenerBaseURL, "/") + "/[a-zA-Z]{8}"

				assert.Regexp(t, regexToValidateTheLink, resp.ShortenedURL)
				assert.Equal(t, testData.want.contentType, res.Header.Get("Content-Type"))
			}

		})
	}

}
