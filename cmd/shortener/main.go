package main

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/compression"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/logging"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	s, err := storage.GetStorage()
	if err != nil {
		panic(err)
	}
	short, _ := app.NewURLShortenService(s)

	// создаём предустановленный регистратор zap
	logger, errLog := zap.NewDevelopment()
	if errLog != nil {
		// вызываем панику, если ошибка
		panic(errLog)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugaredLogger := *logger.Sugar()

	r := mux.NewRouter()
	r.HandleFunc("/", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.RootPageHandler(short))))
	r.HandleFunc("/ping", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.PingHandler(short))))
	r.HandleFunc("/{slug:[A-Za-z]+}", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.ExpandHandler(short))))
	r.HandleFunc("/api/shorten", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.RootPageJSONHandler(short))))
	r.HandleFunc("/api/shorten/batch", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.RootPageJSONBatchHandler(short))))

	err = http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
