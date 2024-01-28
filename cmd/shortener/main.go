package main

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {

	// создаём предустановленный регистратор zap
	logger, errLog := zap.NewDevelopment()
	if errLog != nil {
		// вызываем панику, если ошибка
		panic(errLog)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugaredLogger := *logger.Sugar()

	short := app.NewURLShortenService()

	r := mux.NewRouter()
	r.HandleFunc("/", logging.WithLogging(sugaredLogger, handlers.RootPageHandler(short)))
	r.HandleFunc("/{slug:[A-Za-z]+}", logging.WithLogging(sugaredLogger, handlers.ExpandHandler(short)))

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
