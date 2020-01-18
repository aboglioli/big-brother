package models

type Employee struct {
	Base
	Timestamp
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	RoleID         string `json:"role_id"`
}

func NewEmployee() *Employee {
	return &Employee{
		Base:      NewBase(),
		Timestamp: NewTimestamp(),
	}
}
