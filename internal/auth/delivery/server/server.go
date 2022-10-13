// Package server contains gRPC API endpoint to work with GoK server side service.
package server

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"

	usrUC "github.com/sergeysynergy/gok/internal/auth/useCase/user"
	"github.com/sergeysynergy/gok/internal/entity"
	serviceErrors "github.com/sergeysynergy/gok/internal/errors"
	pb "github.com/sergeysynergy/gok/proto"
)

type AuthServer struct {
	pb.UnimplementedAuthServer
	lg   *zap.Logger
	user usrUC.UseCase
}

func New(
	logger *zap.Logger,
	user usrUC.UseCase,
) *AuthServer {
	return &AuthServer{
		lg:   logger,
		user: user,
	}
}

func (s AuthServer) SignIn(ctx context.Context, in *pb.SignInRequest) (*pb.SignInResponse, error) {
	usr := &entity.User{
		Login: in.User.Login,
	}

	signedUsr, err := s.user.SignIn(ctx, usr)
	if err != nil {
		if errors.Is(err, serviceErrors.ErrUserAlreadyExists) {
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
