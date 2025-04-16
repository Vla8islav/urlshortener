package main

import (
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootPageHandler)
	r.HandleFunc("/{slug:[A-Za-z]+}", handlers.ExpandHandler)

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
