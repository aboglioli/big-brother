package auth

import (
	"time"

	"github.com/aboglioli/big-brother/config"
	"github.com/aboglioli/big-brother/errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrEncode        = errors.Internal.New("auth.token.encode")
	ErrSigningMethod = errors.Internal.New("auth.token.signing_method")
	ErrDecode        = errors.Internal.New("auth.token.decode")
)

type Token struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	UserID    string             `json:"userId" bson:"userId"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

func NewToken(userID string) *Token {
	return &Token{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		CreatedAt: time.Now(),
	}
}

func (t *Token) Encode() (string, error) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        t.ID.Hex(),
		"userId":    t.UserID,
		"createdAt": t.CreatedAt,
	})

	config := config.Get()
	tokenStr, err := jwtToken.SignedString(config.AuthHMACSecret)

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
		return config.AuthHMACSecret, nil
	})

	if err != nil {
		return nil, ErrDecode.Wrap(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(primitive.ObjectID)
		userID := claims["userId"].(string)
		createdAt := claims["createdAt"].(time.Time)
		return &Token{
			ID:        id,
			UserID:    userID,
			CreatedAt: createdAt,
		}, nil
	}

	return nil, ErrDecode
}
