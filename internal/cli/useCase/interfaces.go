package useCase

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Client interface {
	SignIn(context.Context, *entity.User) (*entity.SignedUser, error)
	Init(ctx context.Context, token string) (*entity.Branch, error)
}

type UseCase interface {
	SignIn(*entity.CLIUser) (*entity.SignedUser, error)
	Init(token string) (*entity.Branch, error)
}
