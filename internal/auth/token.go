package auth

import (
	"time"

	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/db"
	"github.com/aboglioli/big-brother/pkg/errors"
	"github.com/dgrijalva/jwt-go"
)

var (
	ErrEncode        = errors.Internal.New("auth.token.encode")
	ErrSigningMethod = errors.Internal.New("auth.token.signing_method")
	ErrDecode        = errors.Internal.New("auth.token.decode")
)

type Token struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

func NewToken(userID string) *Token {
	return &Token{
		ID:        db.NewID(),
		UserID:    userID,
		CreatedAt: time.Now().UnixNano(),
	}
}

func (t *Token) Encode() (string, error) {
	config := config.Get()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        t.ID,
		"userId":    t.UserID,
		"createdAt": t.CreatedAt,
	})

	tokenStr, err := jwtToken.SignedString(config.JWTSecret)
	if err != nil {
		return "", ErrEncode.Wrap(err)
	}

	return tokenStr, nil
}

func decodeToken(tokenStr string) (*Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrSigningMethod
		}

		config := config.Get()
		return config.JWTSecret, nil
	})
	if err != nil {
		return nil, ErrDecode.Wrap(err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrDecode
	}

	id := claims["id"].(string)
	userID := claims["userId"].(string)
	createdAt := int64(claims["createdAt"].(float64))

	return &Token{
		ID:        id,
		UserID:    userID,
		CreatedAt: createdAt,
	}, nil
}
