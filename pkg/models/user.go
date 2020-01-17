package models

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/pkg/errors"
)

var (
	ErrSetPassword = errors.Internal.New("user.set_password")
)

type Role string

const (
	ADMIN = Role("admin")
	USER  = Role("user")
)

type User struct {
	Base
	Username string `json:"username"validate:"required,min=4,max=32,alphanumdash"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,min=5,max=64,email"`
	Name     string `json:"name" validate:"required,min=2,max=32,alphaspaces"`
	Lastname string `json:"lastname" validate:"required,min=2,max=32,alphaspaces"`
	Role     Role   `json:"role"`

	Validated bool `json:"validated" bson:"validated"`
}

func NewUser() *User {
	return &User{
		Base: NewBase(),
		Role: USER,
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
