package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthRequired = status.Error(codes.Unauthenticated, "authentication required")
	ErrPullUpToDate = status.Error(codes.NotFound, "nothing to pull: already up to date")
)
