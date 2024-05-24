package main

import (
	"github.com/Vla8islav/urlshortener/internal/app/api"
	"github.com/Vla8islav/urlshortener/internal/app/compression"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/cookies"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/helpers"
	"github.com/Vla8islav/urlshortener/internal/app/logging"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	ctx, cancel := helpers.GetDefaultContext()
	defer cancel()

	s, err := storage.GetStorage(ctx)
	if err != nil {
		panic(err)
	}
	defer s.Close()

	short, _ := api.NewGoBooruService(ctx, s)

	// создаём предустановленный регистратор zap
	logger, errLog := zap.NewDevelopment()
	if errLog != nil { /* вызываем панику, если ошибка */
		panic(errLog)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugaredLogger := *logger.Sugar()

	r := mux.NewRouter()

	r.HandleFunc("/ping",
		logging.WithLogging(sugaredLogger,
			cookies.SetUserCookie(&s,
				handlers.PingHandler(&s))))

	r.HandleFunc("/register",
		compression.GzipHandle(
			logging.WithLogging(sugaredLogger,
				cookies.SetUserCookie(&s,
					handlers.RegisterHandler(short)))))

	err = http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {

		panic(err)
	}
}
