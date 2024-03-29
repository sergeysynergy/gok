// Package client implements gRPC client for working with GoK storage API.
package client

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	pb "github.com/sergeysynergy/gok/proto"
	"github.com/sergeysynergy/gok/tool/serializers"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"
)

type GokClient struct {
	lg   *zap.Logger
	once *sync.Once
	// gRPC Auth service address.
	authAddr string
	// gRPC Storage service address.
	storageAddr string
}

func New(logger *zap.Logger, authAddr, storageAddr string) *GokClient {
	c := &GokClient{
		lg:          logger,
		authAddr:    authAddr,
		storageAddr: storageAddr,
	}

	return c
}

func (c *GokClient) getAuthConnect() (pb.AuthClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(c.authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewAuthClient(conn), conn
}

func (c *GokClient) getStorageConnect() (pb.StorageClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(c.storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewStorageClient(conn), conn
}

func (c *GokClient) SignIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
	auth, conn := c.getAuthConnect()
	defer conn.Close()

	resp, err := auth.SignIn(ctx, &pb.SignInRequest{
		User: &pb.UserForAdd{
			Login: usr.Login,
		},
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.AlreadyExists {
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrUserAlreadyExists)
			} else {
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("Failed to parse error: " + err.Error())
		return nil, err
	}

	return &entity.SignedUser{
		Token: resp.User.Token,
		Key:   resp.User.Key,
	}, nil
}

func (c *GokClient) Login(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
	auth, conn := c.getAuthConnect()
	defer conn.Close()

	resp, err := auth.Login(ctx, &pb.LoginRequest{
		User: &pb.UserForAdd{
			Login: usr.Login,
		},
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.NotFound {
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrUserNotFound)
			} else {
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("Failed to parse error: " + err.Error())
		return nil, err
	}

	return &entity.SignedUser{
		Token: resp.User.Token,
		Key:   resp.User.Key,
	}, nil
}

func (c *GokClient) Init(ctx context.Context, token string) (*entity.Branch, error) {
	storage, conn := c.getStorageConnect()
	defer conn.Close()

	// Add token value to metadata using context.
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := storage.InitBranch(ctx, &empty.Empty{})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrAuthRequired)
			case codes.NotFound:
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrRecordNotFound)
			default:
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("GokClient.Init - failed to parse error: " + err.Error())
		return nil, err
	}

	return &entity.Branch{
		ID:   entity.BranchID(resp.Branch.Id),
		Name: resp.Branch.Name,
		Head: resp.Branch.Head,
	}, nil
}

func (c *GokClient) Push(ctx context.Context, token string, brn *entity.Branch, recs []*entity.Record) (*entity.Branch, error) {
	var err error
	logPrefix := "GokClient.Push"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			c.lg.Error(err.Error())
		}
	}()

	storage, conn := c.getStorageConnect()
	defer conn.Close()

	// Add token value to metadata using context.
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	reqsPB := serializers.RecordsEntityToPB(recs)

	resp, err := storage.Push(ctx, &pb.PushRequest{
		Branch: &pb.Branch{
			Id:   uint64(brn.ID),
			Name: brn.Name,
			Head: brn.Head,
		},
		Records: reqsPB,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrAuthRequired)
			} else {
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("GokClient.Push - failed to parse error: " + err.Error())
		return nil, err
	}

	brn = &entity.Branch{
		ID:   entity.BranchID(resp.Branch.Id),
		Name: resp.Branch.Name,
		Head: resp.Branch.Head,
	}

	c.lg.Debug(fmt.Sprintf("%s successful, got branch: ID %d; name %s; head %d", logPrefix, brn.ID, brn.Name, brn.Head))
	return brn, nil
}

func (c *GokClient) Pull(ctx context.Context, token string, brn *entity.Branch) (*entity.Branch, []*entity.Record, error) {
	c.lg.Debug("doing GokDeliveryClient.Pull")

	storage, conn := c.getStorageConnect()
	defer conn.Close()

	// Add token value to metadata using context.
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := storage.Pull(ctx, &pb.PullRequest{
		Branch: &pb.Branch{
			Id:   uint64(brn.ID),
			Name: brn.Name,
			Head: brn.Head,
		},
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.Unauthenticated:
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrAuthRequired)
			default:
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("Failed to parse error: " + err.Error())
		return nil, nil, err
	}

	recs := serializers.RecordsPBToEntity(resp.Records)

	return &entity.Branch{
		ID:   entity.BranchID(resp.Branch.Id),
		Name: resp.Branch.Name,
		Head: resp.Branch.Head,
	}, recs, nil
}
