// Package client implements gRPC client for working with GoK storage API.
package client

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"

	pb "github.com/sergeysynergy/gok/proto"
)

type Client struct {
	lg   *zap.Logger
	once *sync.Once
	// gRPC Auth service address.
	authAddr string
	// gRPC Storage service address.
	storageAddr string
}

func New(logger *zap.Logger, authAddr, storageAddr string) *Client {
	c := &Client{
		lg:          logger,
		authAddr:    authAddr,
		storageAddr: storageAddr,
	}

	return c
}

func (c *Client) getAuthConnect() (pb.AuthClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(c.authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewAuthClient(conn), conn
}

func (c *Client) getStorageConnect() (pb.StorageClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(c.storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewStorageClient(conn), conn
}

func (c *Client) SignIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
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

func (c *Client) Init(ctx context.Context, token string) (*entity.Branch, error) {
	storage, conn := c.getStorageConnect()
	defer conn.Close()

	// Add token value to metadata using context.
	md := metadata.New(map[string]string{"token": token})
	ctx = metadata.NewOutgoingContext(ctx, md)

	resp, err := storage.InitBranch(ctx, &empty.Empty{})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			if e.Code() == codes.Unauthenticated {
				err = fmt.Errorf("%s - %w", e.Message(), gokErrors.ErrAuthRequired)
			} else {
				err = fmt.Errorf("%s - %s", e.Code(), e.Message())
			}
		}
		c.lg.Debug("Failed to parse error: " + err.Error())
		return nil, err
	}

	return &entity.Branch{
		Name: resp.Branch.Name,
		Head: resp.Branch.Head,
	}, nil
}
