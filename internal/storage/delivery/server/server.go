// Package server contains gRPC API endpoints to work with `storage` service.
package server

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"

	brnUC "github.com/sergeysynergy/gok/internal/storage/useCase/branch"
	pb "github.com/sergeysynergy/gok/proto"
)

// StorageServer implements API points to work with `Storage` service using gRPC protocol.
type StorageServer struct {
	pb.UnimplementedStorageServer
	lg     *zap.Logger
	branch brnUC.UseCase
}

func New(
	logger *zap.Logger,
	branch brnUC.UseCase,
) *StorageServer {
	return &StorageServer{
		lg:     logger,
		branch: branch,
	}
}

func (s StorageServer) InitBranch(ctx context.Context, _ *pb.InitBranchRequest) (*pb.InitBranchResponse, error) {
	// Get token value from metadata.
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) == 0 {
			return nil, fmt.Errorf("failed to get token from metadata")
		}
		token = values[0]
	}

	// TODO: replace token with userID based on auth interceptor
	brn, err := s.branch.AddGet(ctx, token)
	if err != nil {
		return nil, err
	}

	return &pb.InitBranchResponse{
		Branch: &pb.Branch{
			Name:   brn.Name,
			Head:   brn.Head,
			UserID: uint32(brn.UserID),
		},
	}, nil
}
