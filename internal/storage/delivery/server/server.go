// Package server contains gRPC API endpoints to work with `storage` service.
package server

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	brnUC "github.com/sergeysynergy/gok/internal/storage/useCase/branch"
	pb "github.com/sergeysynergy/gok/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
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

func (s StorageServer) InitBranch(ctx context.Context, _ *empty.Empty) (*pb.InitBranchResponse, error) {
	// Get token value from metadata.
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) == 0 {
			return &pb.InitBranchResponse{}, ErrAuthRequired
		}
		token = values[0]
	}

	brn, err := s.branch.AddGet(ctx, token)
	if err != nil {
		return &pb.InitBranchResponse{}, err
	}

	return &pb.InitBranchResponse{
		Branch: &pb.Branch{
			Name: brn.Name,
			Head: brn.Head,
		},
	}, nil
}
