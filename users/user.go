package users

import (
	"encoding/json"

	"github.com/aboglioli/big-brother/db/models"
	"github.com/aboglioli/big-brother/errors"
	"golang.org/x/crypto/bcrypt"
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
	models.Base
	Username string `json:"username" bson:"username" validate:"required,min=4,max=32,alphanumdash"`
	Password string `json:"password" bson:"password" validate:"required"`
	Email    string `json:"email" bson:"email" validate:"required,email"`
	Name     string `json:"name" bson:"name" validate:"required,min=2,max=32,alphaspaces"`
	Lastname string `json:"lastname" bson:"lastname" validate:"required,min=2,max=32,alphaspaces"`
	Roles    []Role `json:"roles" bson:"roles"`

	Validated bool `json:"validated" bson:"validated"`
}

func NewUser() *User {
	return &User{
		Base:  models.NewBase(),
		Roles: []Role{USER},
	}
}

func (u *User) SetPassword(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return ErrSetPassword.M("cannot generate hash from password %s", pwd).C("password", pwd).Wrap(err)
	}
	u.Password = string(hash)
	return nil
}

func (u *User) ComparePassword(pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)); err != nil {
		return false
	}
	return true
}

func (u *User) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}
