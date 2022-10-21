package useCase

import (
	"context"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	Create(context.Context, *entity.Record) error
	Read(context.Context, entity.RecordID) (*entity.Record, error)
	Update(context.Context, *entity.Record) error

	// TypeList return records of given type.
	TypeList(context.Context, gokConsts.RecordType) ([]*entity.Record, error)
	// HeadList return all records where record head more than given head.
	HeadList(ctx context.Context, head uint64) ([]*entity.Record, error)
	// BulkCreateUpdate records.
	BulkCreateUpdate(context.Context, []*entity.Record) error
}

type Client interface {
	SignIn(context.Context, *entity.User) (*entity.SignedUser, error)
	Login(context.Context, *entity.User) (*entity.SignedUser, error)
	Init(ctx context.Context, token string) (*entity.Branch, error)
	Push(ctx context.Context, token string, branch *entity.Branch, records []*entity.Record) (*entity.Branch, error)
	Pull(ctx context.Context, token string, branch *entity.Branch) (*entity.Branch, []*entity.Record, error)
}

type UseCase interface {
	SignIn(*entity.CLIUser) (*entity.SignedUser, error)
	Login(*entity.CLIUser) (*entity.SignedUser, error)
	Init(token string) (*entity.Branch, error)
	Push(token string, branch string, head uint64) (*entity.Branch, error)
	Pull(token string, branch string, head uint64) (*entity.Branch, error)

	DescAdd(*entity.Record) error
	DescSet(*entity.Record) error
	DescList() ([]*entity.Record, error)
}
