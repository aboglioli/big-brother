package models

type CRUD struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type Permission struct {
	Module Module `json:"module"`
	CRUD   CRUD   `json:"crud"`
}
