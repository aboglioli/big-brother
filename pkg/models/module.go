package models

type Module struct {
	Slug string `json:"slug" validate:"required"`
	Name string `json:"name" validate:"required"`
}
