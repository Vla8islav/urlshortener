package handlers

import "github.com/gorilla/mux"

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", RootPageHandler)
	r.HandleFunc("/{slug:[A-Za-z]+}", ExpandHandler)
	return r
}
