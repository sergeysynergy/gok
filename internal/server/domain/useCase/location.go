package useCase

import (
	"go.uber.org/zap"
)

type LocationUseCase struct {
	lg   *zap.Logger
	repo Repo
}

var _ UseCase = new(LocationUseCase)

func New(logger *zap.Logger, repo Repo) *LocationUseCase {
	s := &LocationUseCase{
		lg:   logger,
		repo: repo,
	}

	return s
}
