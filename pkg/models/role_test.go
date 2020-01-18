package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRolePermissions(t *testing.T) {
	prodMod := Module{
		Slug: "product",
		Name: "Productos",
	}
	orgMod := Module{
		Slug: "organization",
		Name: "Negocio",
	}
	sellMod := Module{
		Slug: "sell",
		Name: "Ventas",
	}
	emplMod := Module{
		Slug: "employee",
		Name: "Empleados",
	}
	msgMod := Module{
		Slug: "message",
		Name: "Mensajes",
	}

	role := Role{
		Name: "Encargado",
		Permissions: []Permission{{
			Module: prodMod,
			CRUD: CRUD{
				Create: true,
				Read:   true,
				Update: true,
			},
		}, {
			Module: orgMod,
			CRUD: CRUD{
				Read:   true,
				Update: true,
			},
		}, {
			Module: sellMod,
			CRUD: CRUD{
				Read: true,
			},
		}, {
			Module: emplMod,
			CRUD: CRUD{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			},
		}, {
			Module: msgMod,
			CRUD: CRUD{
				Read:   true,
				Delete: true,
			},
		}},
	}

	tests := []struct {
		name      string
		role      Role
		module    string
		canCreate bool
		canRead   bool
		canUpdate bool
		canDelete bool
	}{{
		"product module",
		role,
		"product",
		true, true, true, false,
	}, {
		"organization module",
		role,
		"organization",
		false, true, true, false,
	}, {
		"sell module",
		role,
		"sell",
		false, true, false, false,
	}, {
		"employee module",
		role,
		"employee",
		true, true, true, true,
	}, {
		"message module",
		role,
		"message",
		false, true, false, true,
	}, {
		"not existing module",
		role,
		"non-existing",
		false, false, false, false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			crud := test.role.Privileges(test.module)
			assert.Equal(test.canCreate, crud.Create)
			assert.Equal(test.canRead, crud.Read)
			assert.Equal(test.canUpdate, crud.Update)
			assert.Equal(test.canDelete, crud.Delete)
		})
	}
}
