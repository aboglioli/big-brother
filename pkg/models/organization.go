package models

type Organization struct {
	Base
	Timestamp
	Name string `json:"name" validate:"required,min=4,max=64"`
}

func NewOrganization() *Organization {
	return &Organization{
		Base:      NewBase(),
		Timestamp: NewTimestamp(),
	}
}

func (o *Organization) Clone() *Organization {
	c := *o
	return &c
}
