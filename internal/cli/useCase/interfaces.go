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
	TypeList(context.Context, entity.BranchID, gokConsts.RecordType) ([]*entity.Record, error)
	// HeadList return all records where record head more than given head.
	HeadList(ctx context.Context, brnID entity.BranchID, head uint64) ([]*entity.Record, error)
	// BulkCreateUpdate records.
	BulkCreateUpdate(context.Context, []*entity.Record) error
	// ByIDsList return all records from given IDs slice.
	ByIDsList(ctx context.Context, ids []string) ([]*entity.Record, error)
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
	Init(token string, head uint64) (*entity.Branch, error)
	Push(token string, branch *entity.Branch) (*entity.Branch, error)
	Pull(*entity.CLIConf, *entity.Branch) (*entity.Branch, error)

	RecordAdd(*entity.CLIConf, *entity.Record) error
	RecordSet(*entity.CLIConf, *entity.Record) error
	RecordList(conf *entity.CLIConf, recType gokConsts.RecordType) ([]*entity.Record, error)
}
