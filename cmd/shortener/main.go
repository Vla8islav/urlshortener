package main

import (
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	short := app.URLShorten{S: storage.GetMakeshiftStorageInstance()}
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.RootPageHandler(&short))
	r.HandleFunc("/{slug:[A-Za-z]+}", handlers.ExpandHandler(&short))

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, r)
	if err != nil {
		panic(err)
	}
}
