package handlers

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExpandHandler(t *testing.T) {
	type expectedResult struct {
		code int
	}

	validRequest := httptest.NewRequest(http.MethodPost, "/", nil)
	validRequest.Header = http.Header{
		"Content-Type": []string{"text/plain"},
	}
	validRequest.Body = io.NopCloser(strings.NewReader("http://ya.ru"))

	testDataArray := []struct {
		name    string
		request *http.Request
		want    expectedResult
	}{
		{
			name:    "Successful link generation",
			request: validRequest,
			want:    expectedResult{code: 201},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			ExpandHandler(w, testData.request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
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
