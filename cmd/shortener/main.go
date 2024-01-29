package main

import (
	"compress/gzip"
	"github.com/Vla8islav/urlshortener/internal/app"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/handlers"
	"github.com/Vla8islav/urlshortener/internal/app/logging"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		// это упрощённый пример. В реальном приложении следует проверять все
		// значения r.Header.Values("Accept-Encoding") и разбирать строку
		// на составные части, чтобы избежать неожиданных результатов
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func main() {

	// создаём предустановленный регистратор zap
	logger, errLog := zap.NewDevelopment()
	if errLog != nil {
		// вызываем панику, если ошибка
		panic(errLog)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugaredLogger := *logger.Sugar()

	short := app.NewURLShortenService()

	r := mux.NewRouter()
	r.HandleFunc("/", logging.WithLogging(sugaredLogger, handlers.RootPageHandler(short)))
	r.HandleFunc("/{slug:[A-Za-z]+}", logging.WithLogging(sugaredLogger, handlers.ExpandHandler(short)))
	r.HandleFunc("/api/shorten", logging.WithLogging(sugaredLogger, handlers.RootPageJSONHandler(short)))

	err := http.ListenAndServe(configuration.ReadFlags().ServerAddress, gzipHandle(r))
	if err != nil {

		panic(err)
	}
}
