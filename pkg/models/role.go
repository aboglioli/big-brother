package models

type Role struct {
	Base
	Name        string       `json:"name" validate:"required,min=3,max=32"`
	Permissions []Permission `json:"permissions"`
}

func NewRole() Role {
	return Role{
		Base: NewBase(),
	}
}

func (r Role) Privileges(m string) CRUD {
	for _, p := range r.Permissions {
		if p.Module.Slug == m {
			return p.CRUD
		}
	}
	return CRUD{}
}
