package models

import (
	"time"
)

type Token struct {
	Base
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

func NewToken(userID string) *Token {
	return &Token{
		Base:      NewBase(),
		UserID:    userID,
		CreatedAt: time.Now().UnixNano(),
	}
}
