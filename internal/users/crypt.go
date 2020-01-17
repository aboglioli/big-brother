package users

import (
	"github.com/aboglioli/big-brother/pkg/config"
	"github.com/aboglioli/big-brother/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrCryptHash = errors.Internal.New("crypt.hash")
)

// Interface
type PasswordCrypt interface {
	Hash(pwd string) (string, error)
	Compare(hashedPwd, pwd string) bool
}

// Implementation
type bcryptCrypt struct {
	cost int
}

func NewBcryptCrypt() PasswordCrypt {
	c := config.Get()
	return &bcryptCrypt{
		cost: c.BcryptCost,
	}
}

func (b *bcryptCrypt) Hash(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), b.cost)
	if err != nil {
		return "", ErrCryptHash.M("cannot generate hash from password %s", pwd).C("password", pwd).Wrap(err)
	}
	return string(hash), nil
}

func (b *bcryptCrypt) Compare(hashedPwd, pwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd)); err != nil {
		return false
	}
	return true
}
