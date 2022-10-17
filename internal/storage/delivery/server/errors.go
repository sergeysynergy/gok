package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserUnknownError    = status.Error(codes.Unknown, "user unknown error")
	ErrUserInvalidArgument = status.Error(codes.InvalidArgument, "invalid argument for user")
	ErrUserAlreadyExists   = status.Error(codes.AlreadyExists, "user already exists")
	ErrUserNotFound        = status.Error(codes.NotFound, "user not found")
	ErrUsersList           = status.Error(codes.Unknown, "failed to get users list")
)
