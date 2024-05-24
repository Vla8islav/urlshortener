package cookies

import (
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/http"
	"net/url"
)

func SetUserCookie(storage *storage.Storage, next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции

		baseURL := configuration.ReadFlags().ShortenerBaseURL
		_, err := url.Parse(baseURL)
		if err != nil {
			panic("Base URL isn't parsable url " + baseURL + err.Error())
		}

		cookieName := "userid"
		_, err = r.Cookie(cookieName)

		// передаём управление хендлеру
		next.ServeHTTP(w, r)

	}

}
