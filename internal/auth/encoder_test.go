package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncoderEncode(t *testing.T) {
	assert := assert.New(t)

	enc := &encoder{[]byte("my_secret")}

	tokenStr, err := enc.Encode("token123")
	assert.Nil(err)
	assert.NotEmpty(tokenStr)
	assert.Greater(len(tokenStr), 10)
}

func TestEncoderDecode(t *testing.T) {
	assert := assert.New(t)

	enc := &encoder{[]byte("my_secret")}

	// Valid secret: my_secret
	tokenID, err := enc.Decode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InRva2VuMTIzIn0.EZljBLKD2KvaRZSQDdHQ0IEAsnK7u5Dx2cdGHdVg0OQ")
	assert.Nil(err)
	assert.Equal(tokenID, "token123")

	// Invalid secret: invalid_secret
	tokenID, err = enc.Decode("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6InRva2VuMTIzIn0.stdAmb1gRntWULyhf9EmFGuhDVrMdFgfRELryRsQtqA")
	assert.NotNil(err)
	assert.Empty(tokenID)
}
