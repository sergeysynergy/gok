// Package grpc contains implementation of gRPC server.
package grpc

import (
	"go.uber.org/zap"

	client "back-mfsb/api/pb_go"
	"back-mfsb/service/location/internal/domain/useCase"
	"back-mfsb/service/neotuner/pkg/protowrapper"
)

// LocationService реализует методы сгенерённого gRPC-сервиса.
type LocationService struct {
	protowrapper.ProtoWrapper
	lg   *zap.Logger
	cc   client.AuthorizationClient
	repo useCase.Repo
	uc   useCase.UseCase
}

// TODO: добавить проверку на реализацию контракта generated.LocationService

func New(
	logger *zap.Logger,
	cc client.AuthorizationClient,
	repo useCase.Repo,
	uc useCase.UseCase,
) *LocationService {
	return &LocationService{
		lg:   logger,
		cc:   cc,
		repo: repo,
		uc:   uc,
	}
}
