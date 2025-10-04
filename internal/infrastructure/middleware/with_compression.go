package middleware

import (
	"io"
	"net/http"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func WithGzipCompression() func(http.Handler) http.Handler {

	return func(handlerNext http.Handler) http.Handler {

		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				//if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				//	gzw := gzip.NewWriter(w)
				//	defer func(gzw *gzip.Writer) {
				//		err := gzw.Close()
				//		if err != nil {
				//			panic(err)
				//		}
				//	}(gzw)
				//	w.Header().Set("Content-Encoding", "gzip")
				//	grw := gzipResponseWriter{ResponseWriter: w, Writer: gzw}
				//	handlerNext.ServeHTTP(grw, r)
				//	return
				//}

				handlerNext.ServeHTTP(w, r)

			},
		)

	}

}
