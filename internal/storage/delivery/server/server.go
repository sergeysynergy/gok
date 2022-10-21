// Package server contains gRPC API endpoints to work with `storage` service.
package server

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	brnUC "github.com/sergeysynergy/gok/internal/storage/useCase/branch"
	recUC "github.com/sergeysynergy/gok/internal/storage/useCase/record"
	pb "github.com/sergeysynergy/gok/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			Name: brn.Name,
			Head: brn.Head,
		},
	}, nil
}

func (s StorageServer) Push(ctx context.Context, in *pb.PushRequest) (*pb.PushResponse, error) {
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
		Name: in.Branch.Name,
		Head: in.Branch.Head,
	}

	records := make([]*entity.Record, 0, len(in.Records))
	for _, v := range in.Records {
		records = append(records, &entity.Record{
			ID:          entity.RecordID(v.Id),
			Head:        v.Head,
			Branch:      v.Branch,
			Description: entity.Description(v.Description),
			Type:        gokConsts.RecordType(v.Type),
			UpdatedAt:   v.UpdatedAt.AsTime(),
		})
	}

	brn, err := s.branch.Push(ctx, token, brn, records)
	if err != nil {
		return nil, err
	}

	return &pb.PushResponse{
		Branch: &pb.Branch{
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
			Name: in.Branch.Name,
			Head: in.Branch.Head,
		},
	)
	if err != nil {
		if errors.Is(err, gokErrors.ErrPullUpToDate) {
			return nil, ErrPullUpToDate
		}
		return nil, err
	}

	recsPB := make([]*pb.Record, 0, len(recs))
	for _, v := range recs {
		recsPB = append(recsPB, &pb.Record{
			Id:          string(v.ID),
			Head:        v.Head,
			Branch:      v.Branch,
			Description: string(v.Description),
			Type:        string(v.Type),
			UpdatedAt:   timestamppb.New(v.UpdatedAt),
		})
	}

	return &pb.PullResponse{
		Branch: &pb.Branch{
			Name: freshBrn.Name,
			Head: freshBrn.Head,
		},
		Records: recsPB,
	}, nil
}
