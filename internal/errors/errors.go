// Package errors provides description for common project errors.
package errors

import (
	"errors"
)

var (
	ErrAuthRequired = errors.New("authentication required")

	ErrUserUnknown         = errors.New("user unknown error")
	ErrUserInvalidArgument = errors.New("invalid argument for user")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserZeroID          = errors.New("got zero user id")

	ErrSessionNotFound = errors.New("session not found")

	ErrRecordUnknown       = errors.New("record unknown error")
	ErrRecordNotFound      = errors.New("record not found")
	ErrRecordAlreadyExists = errors.New("record already exists")

	ErrPushUnknown = errors.New("push unknown error")

	ErrLocalBranchBehind = errors.New("local branch in behind server")
)
