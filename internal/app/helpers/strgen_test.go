package helpers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateString(t *testing.T) {

	tests := []struct {
		testName     string
		len          int
		letterString string
	}{
		{
			testName:     "Negative number",
			len:          len("EwHXdJfB"),
			letterString: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		},
	}

	for _, testData := range tests {
		t.Run(testData.testName, func(t *testing.T) {
			randomString := GenerateString(testData.len, testData.letterString)
			assert.Regexp(t, "[A-Za-z]+", randomString)
			assert.Len(t, randomString, testData.len)
		})

	}
}
