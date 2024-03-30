package auth

import (
	"fmt"
	"github.com/Vla8islav/urlshortener/internal/app/configuration"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

const DefaultUserID = 1

//func main() {
//	tokenString, err := BuildJWTString()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Println(GetUserID(tokenString))
//}

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const TokenExp = time.Hour * 24 * 365

func GetBearerFromBearerHeader(bearerHeader string) string {
	return strings.Replace(bearerHeader, "Bearer ", "", 1)
}

func GetUserID(tokenString string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(configuration.ReadFlags().SecretKey), nil
		})
	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(userID int) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		// собственное утверждение
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(configuration.ReadFlags().SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}
