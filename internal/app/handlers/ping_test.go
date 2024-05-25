package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/app/api"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPingHandler(t *testing.T) {
	ctx, cancel := helpers.GetDefaultContext()
	defer cancel()

	s, err := storage.NewPostgresStorage(ctx)

	if err != nil {
		panic(err)
	}
	gobooru, _ := api.NewGoBooruService(ctx, s)

	type expectedResult struct {
		code int
	}

	testDataArray := []struct {
		name    string
		request func() *http.Request
		want    expectedResult
	}{
		{
			name: "Successful ping",
			request: func() *http.Request {
				if err != nil {
					panic("Couldn't do a ping test")
				}

				validRequest := httptest.NewRequest(http.MethodGet, "/ping/", nil)
				return validRequest

			},
			want: expectedResult{code: http.StatusOK},
		},
	}

	for _, testData := range testDataArray {
		t.Run(testData.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()

			PingHandler(&gobooru.Storage)(w, testData.request())

			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, testData.want.code, res.StatusCode)

		})
	}

}
