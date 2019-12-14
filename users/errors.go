package users

import (
	"github.com/aboglioli/big-brother/errors"
)

var (
	ErrSetPassword = errors.Internal.New("user.set_password")
)
