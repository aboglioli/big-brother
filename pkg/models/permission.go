package models

type Permission string

const (
	CREATE = Permission("c")
	READ   = Permission("r")
	UPDATE = Permission("u")
	DELETE = Permission("d")
)
