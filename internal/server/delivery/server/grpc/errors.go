package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrLocationUnknownError = status.Error(codes.Unknown, "location unknown error")
	ErrLocationInvalid      = status.Error(codes.InvalidArgument, "invalid argument for location")
	ErrLocationNotFound     = status.Error(codes.NotFound, "location not found")

	ErrLocationsList = status.Error(codes.Unknown, "failed to get locations list")
)
