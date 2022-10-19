package useCase

import (
	"context"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	Create(context.Context, *entity.Record) error
	Update(context.Context, *entity.Record) error
	Read(context.Context, entity.RecordID) (*entity.Record, error)
	List(context.Context, gokConsts.RecordType) ([]*entity.Record, error)
	ListForPush(context.Context, uint64) ([]*entity.Record, error)
}

type Client interface {
	SignIn(context.Context, *entity.User) (*entity.SignedUser, error)
	Login(context.Context, *entity.User) (*entity.SignedUser, error)
	Init(ctx context.Context, token string) (*entity.Branch, error)
	Push(ctx context.Context, token string, branch *entity.Branch, list []*entity.Record) (*entity.Branch, error)
}

type UseCase interface {
	SignIn(*entity.CLIUser) (*entity.SignedUser, error)
	Login(*entity.CLIUser) (*entity.SignedUser, error)
	Init(token string) (*entity.Branch, error)
	Push(token string, branch string, head uint64) (*entity.Branch, error)

	DescAdd(*entity.Record) error
	DescSet(*entity.Record) error
	DescList() ([]*entity.Record, error)
}
