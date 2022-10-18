package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAuthRequired = status.Error(codes.Unauthenticated, "authentication required")
)
