package models

import (
	"time"
)

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

func NewUserDTO(user *User) *UserDTO {
	return &UserDTO{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
		Lastname: user.Lastname,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Validated: user.Validated,
	}
}
