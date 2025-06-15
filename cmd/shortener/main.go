package main

import (
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"net/http"
)

func main() {
	r := handlers.InitRouter()

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
