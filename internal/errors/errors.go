// Package errors provides description for common project errors.
package errors

import (
	"errors"
)

var (
	ErrUserUnknown         = errors.New("user unknown error")
	ErrUserInvalidArgument = errors.New("invalid argument for user")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
)
