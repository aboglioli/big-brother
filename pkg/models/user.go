package models

import (
	"encoding/json"
)

type User struct {
	Base
	Timestamp
	Username string `json:"username" validate:"required,min=4,max=32,alphanumdash"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,min=5,max=64,email"`
	Name     string `json:"name" validate:"required,min=2,max=32,alphaspaces"`
	Lastname string `json:"lastname" validate:"required,min=2,max=32,alphaspaces"`
	Role     string `json:"role"`

	Validated bool `json:"validated" bson:"validated"`
}

func NewUser() *User {
	return &User{
		Base:      NewBase(),
		Timestamp: NewTimestamp(),
		Role:      "user",
	}
}

func (u *User) Clone() *User {
	c := *u
	return &c
}

func (u *User) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}
