package models

import (
	"time"
)

type Token struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt int64  `json:"created_at"`
}

func NewToken(userID string) *Token {
	return &Token{
		ID:        NewID(),
		UserID:    userID,
		CreatedAt: time.Now().UnixNano(),
	}
}
