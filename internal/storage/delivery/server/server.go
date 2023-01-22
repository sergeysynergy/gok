// Package server contains gRPC API endpoints to work with `storage` service.
package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	brnUC "github.com/sergeysynergy/gok/internal/storage/useCase/branch"
	recUC "github.com/sergeysynergy/gok/internal/storage/useCase/record"
	pb "github.com/sergeysynergy/gok/proto"
	"github.com/sergeysynergy/gok/tool/serializers"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

// StorageServer implements API points to work with `Storage` service using gRPC protocol.
type StorageServer struct {
	pb.UnimplementedStorageServer
	lg     *zap.Logger
	branch brnUC.UseCase
	record recUC.UseCase
}

func New(
	logger *zap.Logger,
	branch brnUC.UseCase,
	record recUC.UseCase,
) *StorageServer {
	return &StorageServer{
		lg:     logger,
		branch: branch,
		record: record,
	}
}

func (s StorageServer) InitBranch(ctx context.Context, _ *empty.Empty) (*pb.InitBranchResponse, error) {
	// Get token value from metadata.
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) == 0 {
			return nil, ErrAuthRequired
		}
		token = values[0]
	}

	brn, err := s.branch.AddGet(ctx, token)
	if err != nil {
		return nil, err
	}

	return &pb.InitBranchResponse{
		Branch: &pb.Branch{
			Id:   uint64(brn.ID),
			Name: brn.Name,
			Head: brn.Head,
		},
	}, nil
}

func (s StorageServer) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error) {
	logPrefix := "StorageServer.Push"
	// Get token value from metadata.
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) == 0 {
			return nil, ErrAuthRequired
		}
		token = values[0]
	}

	brn := &entity.Branch{
		ID:   entity.BranchID(in.Branch.Id),
		Name: in.Branch.Name,
		Head: in.Branch.Head,
	}

	recs := serializers.RecordsPBToEntity(in.Records)

	brn, err := s.branch.Push(ctx, token, brn, recs)
	if err != nil {
		return nil, err
	}

	s.lg.Debug(fmt.Sprintf("%s successful, got branch: ID %d; name %s; head %d", logPrefix, brn.ID, brn.Name, brn.Head))
	return &pb.PushResponse{
		Branch: &pb.Branch{
			Id:   uint64(brn.ID),
			Name: brn.Name,
			Head: brn.Head,
		},
	}, nil
}

func (s StorageServer) Pull(ctx context.Context, in *pb.PullRequest) (*pb.PullResponse, error) {
	s.lg.Debug("doing StorageServer.Pull")
	// Get token value from metadata.
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("token")
		if len(values) == 0 {
			return nil, ErrAuthRequired
		}
		token = values[0]
	}

	freshBrn, recs, err := s.branch.Pull(
		ctx,
		token,
		&entity.Branch{
			ID:   entity.BranchID(in.Branch.Id),
			Name: in.Branch.Name,
			Head: in.Branch.Head,
		},
	)
	if err != nil {
		if errors.Is(err, gokErrors.ErrPullUpToDate) {
			return nil, ErrPullUpToDate
		}
		if errors.Is(err, gokErrors.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	recsPB := serializers.RecordsEntityToPB(recs)

	return &pb.PullResponse{
		Branch: &pb.Branch{
			Id:   uint64(freshBrn.ID),
			Name: freshBrn.Name,
			Head: freshBrn.Head,
		},
		Records: recsPB,
	}, nil
}
