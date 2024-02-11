package handlers

import (
	"context"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/jackc/pgx/v5"
	"net/http"
)

func PingHandler(short app.URLShortenServiceMethods) http.HandlerFunc {
	if short == nil {
		panic("Underlying infrastructure isn't initialised")
	}

	return func(res http.ResponseWriter, req *http.Request) {

		conn, err := pgx.Connect(context.Background(), configuration.ReadFlags().DBConnectionString)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		} else {
			defer conn.Close(context.Background())
			res.WriteHeader(http.StatusOK)
		}

	}

}
