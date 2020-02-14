package auth

import (
	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrTokenEncode        = errors.Internal.New("auth.token.encode")
	ErrTokenSigningMethod = errors.Internal.New("auth.token.signing_method")
	ErrTokenDecode        = errors.Internal.New("auth.token.decode")
)

// Interface
type Encoder interface {
	Encode(tokenID string) (string, error)
	Decode(tokenStr string) (string, error)
}

// Implementation
type encoder struct {
	secret []byte
}

func NewEncoder() *encoder {
	c := config.Get()
	return &encoder{
		secret: c.JWTSecret,
	}
}

func (e *encoder) Encode(tokenID string) (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": tokenID,
	})

	tokenStr, err := jwtToken.SignedString(e.secret)
	if err != nil {
		return "", ErrTokenEncode.Wrap(err)
	}

	return tokenStr, nil
}

func (e *encoder) Decode(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSigningMethod
		}
		return e.secret, nil
	})
	if err != nil {
		return "", ErrTokenDecode.Wrap(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", ErrTokenDecode
	}

	id, ok := claims["id"].(string)
	if !ok {
		return "", ErrTokenDecode
	}

	return id, nil
}
