package models

import (
	"encoding/json"
	"time"
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

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`

	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	Validated bool       `json:"validated"`
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Lastname: u.Lastname,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Validated: u.Validated,
	}
}
