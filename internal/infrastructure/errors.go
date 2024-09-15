package infrastructure

import "errors"

var (
	ErrObjectNotFound   = errors.New("object not found")
	ErrAuthUserNotFound = errors.New("user not found")
	ErrAuthInvalidCred  = errors.New("invalid credentials")
)
