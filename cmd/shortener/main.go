package main

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var sugar zap.SugaredLogger

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()

		// эндпоинт /ping
		uri := r.RequestURI
		// метод запроса
		method := r.Method

		// точка, где выполняется хендлер pingHandler
		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)

		// отправляем сведения о запросе в zap
		sugar.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)

	}
	return logFn
}

func main() {

	// создаём предустановленный регистратор zap
	logger, errLog := zap.NewDevelopment()
	if errLog != nil {
		// вызываем панику, если ошибка
		panic(errLog)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()

	short := app.NewURLShortenService()

	r := mux.NewRouter()
	r.HandleFunc("/", WithLogging(handlers.RootPageHandler(short)))
	r.HandleFunc("/{slug:[A-Za-z]+}", WithLogging(handlers.ExpandHandler(short)))

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
