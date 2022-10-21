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

	ErrPushUnknownError = errors.New("push unknown error")

	ErrLocalBranchBehind = errors.New("local branch in behind server")

	ErrPullFailed       = errors.New("pull failed")
	ErrPullUnknownError = errors.New("pull unknown error")
	ErrPullUpToDate     = errors.New("pull already up to date")

	ErrMergeFailed     = errors.New("merge failed")
	ErrResolveConflict = errors.New("resolving merge conflict failed")
	ErrCloningRecord   = errors.New("cloning record failed")
)
