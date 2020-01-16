package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	assert := assert.New(t)
	// Create
	token := NewToken("1234")
	assert.NotNil(token)

	// Encode
	tokenStr, err := token.Encode()
	assert.Nil(err)
	assert.NotEmpty(tokenStr)

	// Decode
	decodedToken, err := decodeToken(tokenStr)
	assert.Nil(err)
	if assert.NotNil(decodedToken) {
		assert.Equal(token.ID, decodedToken.ID)
		assert.Equal(token.UserID, decodedToken.UserID)
	}
}
