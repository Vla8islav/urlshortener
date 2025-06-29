package main

import (
	"github.com/Vla8islav/urlshortener/internal/application"
	"github.com/Vla8islav/urlshortener/internal/handlers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/db"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/logger"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/middleware"
	"net/http"
)

func main() {

	sugaredLogger := logger.NewSugaredLogger()

	repo := db.GetInstance()
	service := application.NewURLShortenService(repo)
	handler := handlers.NewHandler(service)
	router := handlers.InitRouter(handler)

	router.Use(middleware.WithLogging(sugaredLogger))

	sugaredLogger.Infow(
		"Starting server",
		"addr: ", config.ReadFlags().ServerAddress,
	)
	err := http.ListenAndServe(config.ReadFlags().ServerAddress, router)
	if err != nil {
		sugaredLogger.Fatalw("Server encountered an unrecoverable error",
			"error body:", err.Error())

	}
}
