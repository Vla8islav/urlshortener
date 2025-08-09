package middleware

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func WithLogging(sugaredLogger *zap.SugaredLogger) func(http.Handler) http.Handler {

	return func(handlerNext http.Handler) http.Handler {

		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()

				uri := r.RequestURI
				method := r.Method

				handlerNext.ServeHTTP(w, r)

				duration := time.Since(start)

				sugaredLogger.Infoln(
					"uri", uri,
					"method", method,
					"duration", duration,
				)

			},
		)

	}

}
