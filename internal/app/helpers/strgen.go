package helpers

import (
	"fmt"
	"math/rand"
)

func GenerateString(generatedStringLength int, letterString string) string {
	letterStringLen := len(letterString)
	b := make([]byte, generatedStringLength)
	for i := range b {
		b[i] = letterString[rand.Intn(letterStringLen)]
	}
	return string(b)
}

func GenerateRandomURL() string {
	const URLSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return fmt.Sprintf("http://testurl-%s.com/%s",
		GenerateString(4, URLSymbols),
		GenerateString(4, URLSymbols))
}
