package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthEncryption(t *testing.T) {
	userID := 18

	jwt, err := BuildJWTString(userID)
	if err != nil {
		panic("Couldn't create jwt string " + err.Error())
	}

	t.Run("Test encryption decryption", func(t *testing.T) {
		reconstructedUserID, err := GetUserID(jwt)
		assert.Nil(t, err)
		assert.Equal(t, reconstructedUserID, userID)

	})

}
