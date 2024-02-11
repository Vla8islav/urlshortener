package main

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/compression"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	short, _ := app.NewURLShortenService()

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
	r.HandleFunc("/ping", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.PingHandler(short))))
	r.HandleFunc("/", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.RootPageHandler(short))))
	r.HandleFunc("/{slug:[A-Za-z]+}", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.ExpandHandler(short))))
	r.HandleFunc("/api/shorten", logging.WithLogging(sugaredLogger, compression.GzipHandle(handlers.RootPageJSONHandler(short))))

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
