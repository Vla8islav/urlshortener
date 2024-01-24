package helpers

import "math/rand"

func GenerateString(generatedStringLength int, letterString string) string {
	letterStringLen := len(letterString)
	b := make([]byte, generatedStringLength)
	for i := range b {
		b[i] = letterString[rand.Intn(letterStringLen)]
	}
	return string(b)
}
