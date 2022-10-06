// Package errors Пакет предназначен для описания всех ошибок сервиса.
package errors

import (
	"errors"
)

var (
	ErrLocationUnknown  = errors.New("location unknown error")
	ErrLocationInvalid  = errors.New("invalid argument for location")
	ErrLocationNotFound = errors.New("location not found")
)
