package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	gzw *gzip.Writer
}

func (grw *gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gzw.Write(b)
}
func (grw *gzipResponseWriter) WriteHeader(statusCode int) {
	if statusCode != http.StatusNoContent && statusCode != http.StatusNotModified {
		grw.Header().Del("Content-Length")
		grw.Header().Set("Content-Encoding", "gzip")
	}
	grw.ResponseWriter.WriteHeader(statusCode)
}

func (grw *gzipResponseWriter) Flush() {
	err := grw.gzw.Flush()
	if err != nil {
		return
	}
	if flusher, ok := grw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func WithGzipCompression() func(http.Handler) http.Handler {

	return func(handlerNext http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// unpacking
				if strings.EqualFold(r.Header.Get("Accept-Encoding"), "gzip") {
					gr, err := gzip.NewReader(r.Body)
					if err != nil {
						http.Error(w, "failed to use gzip to read a body", http.StatusInternalServerError)
						return
					}
					defer gr.Close()

					r.Body = io.NopCloser(gr)
					r.ContentLength = -1

					r.Header.Del("Content-Encoding")
				}

				//packing
				acceptsGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
				if !acceptsGzip {
					handlerNext.ServeHTTP(w, r)
					return
				}

				if w.Header().Get("Content-Encoding") != "" {
					handlerNext.ServeHTTP(w, r)
					return
				}

				gzw := gzip.NewWriter(w)
				defer gzw.Close()
				grw := &gzipResponseWriter{ResponseWriter: w, gzw: gzw}
				handlerNext.ServeHTTP(grw, r)

			},
		)

	}

}
