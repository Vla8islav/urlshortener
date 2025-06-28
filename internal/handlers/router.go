package handlers

import (
	"github.com/Vla8islav/urlshortener/internal/domain"
	"github.com/gorilla/mux"
)

func InitRouter(handler *Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.RootPageHandler)
	r.HandleFunc("/{slug:[A-Za-z]+}", handler.ExpandHandler)
	return r
}

type Handler struct {
	Service domain.URLShortenService
}

func NewHandler(service domain.URLShortenService) *Handler {
	return &Handler{Service: service}
}
