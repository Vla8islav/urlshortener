package main

import (
	"github.com/Vla8islav/urlshortener/internal/application"
	"github.com/Vla8islav/urlshortener/internal/handlers"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/config"
	"github.com/Vla8islav/urlshortener/internal/infrastructure/db"
	"net/http"
)

func main() {

	repo := db.GetInstance()
	service := application.NewURLShortenService(repo)
	handler := handlers.NewHandler(service)

	r := handlers.InitRouter(handler)

	err := http.ListenAndServe(config.ReadFlags().ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
