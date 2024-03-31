package cookies

import (
	"errors"
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/auth"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/Vla8islav/urlshortener/internal/app/storage"
	"net/http"
	"net/url"
	"time"
)

func SetUserCookie(storage *storage.Storage, next http.Handler) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции

		expire := time.Now().AddDate(1, 0, 0)
		baseURL := configuration.ReadFlags().ShortenerBaseURL
		u, err := url.Parse(baseURL)
		if err != nil {
			panic("Base URL isn't parsable url " + baseURL + err.Error())
		}

		var userID int
		cookieName := "userid"
		_, err = r.Cookie(cookieName)
		if errors.Is(err, http.ErrNoCookie) {
			userID, err = (*storage).GetNewUserID(r.Context())
			if err != nil {
				panic("Couldn't create a new user" + err.Error())
			}

			cookieValue, err := auth.BuildJWTString(userID)
			if err != nil {
				panic(fmt.Sprintf("Couldn't build a token string out of user %d %s", userID, err.Error()))
			}

			cookie := http.Cookie{Name: cookieName,
				Value:      cookieValue,
				Path:       "/",
				Domain:     u.Hostname(),
				Expires:    expire,
				RawExpires: expire.Format(time.UnixDate),
				MaxAge:     365 * 24 * 60 * 60,
				Secure:     true,
				HttpOnly:   true,
				SameSite:   http.SameSiteDefaultMode,
				Raw:        cookieName + "=" + cookieValue,
				Unparsed:   []string{cookieName + "=" + cookieValue}}
			http.SetCookie(w, &cookie)
			r.AddCookie(&cookie)
		}

		// передаём управление хендлеру
		next.ServeHTTP(w, r)

	}

}
