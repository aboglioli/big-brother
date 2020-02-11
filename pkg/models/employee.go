package models

type Employee struct {
	Base
	Timestamp
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	Role           Role   `json:"role"`
}

func NewEmployee() *Employee {
	return &Employee{
		Base:      NewBase(),
		Timestamp: NewTimestamp(),
	}
}
