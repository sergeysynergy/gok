// Package server contains gRPC API endpoint to work with GoK server side service.
package server

import (
	"context"
	"errors"
	"fmt"
	usrUC "github.com/sergeysynergy/gok/internal/auth/useCase/user"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	pb "github.com/sergeysynergy/gok/proto"
	"go.uber.org/zap"
	"net"
)

// AuthServer implements API points to work with `Auth` service using gRPC protocol.
type AuthServer struct {
	pb.UnimplementedAuthServer
	lg            *zap.Logger
	trustedSubnet *net.IPNet
	user          usrUC.UseCase
}

func New(
	logger *zap.Logger,
	trustedSubnet *net.IPNet,
	user usrUC.UseCase,
) *AuthServer {
	return &AuthServer{
		lg:            logger,
		trustedSubnet: trustedSubnet,
		user:          user,
	}
}

// SignIn new user: create new user record, create new session, return session token.
func (s AuthServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	usr := &entity.User{
		Login: in.User.Login,
	}

	signedUsr, err := s.user.SignIn(ctx, usr)
	if err != nil {
		if errors.Is(err, gokErrors.ErrUserAlreadyExists) {
			return nil, ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("%w - %s", ErrUserUnknownError, err)
	}

	return &pb.SignInResponse{
		User: &pb.SignedUser{
			Token: signedUsr.Token,
			Key:   signedUsr.Key,
		},
	}, nil
}

// Login already signed users.
func (s AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	usr := &entity.User{
		Login: in.User.Login,
	}
	fmt.Println("::", in.User)

	signedUsr, err := s.user.Login(ctx, usr)
	if err != nil {
		if errors.Is(err, gokErrors.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w - %s", ErrUserUnknownError, err)
	}

	return &pb.LoginResponse{
		User: &pb.SignedUser{
			Token: signedUsr.Token,
			Key:   signedUsr.Key,
		},
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Internal methods for cross-service communication.

func (s AuthServer) cidrCheck() error {
	if s.trustedSubnet.IP == nil || s.trustedSubnet.Mask == nil {
		return fmt.Errorf("invalid trusted subnet given")
	}

	// TODO: Add CIDR check

	return nil
}

// GetUser by token.
func (s AuthServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	err := s.cidrCheck()
	if err != nil {
		return nil, ErrAuthRequired
	}

	usr, err := s.user.Get(ctx, in.Token)
	if err != nil {
		if errors.Is(err, gokErrors.ErrSessionNotFound) {
			return nil, ErrAuthRequired
		}
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			ID:    int32(usr.ID),
			Login: usr.Login,
		},
	}, nil
}
