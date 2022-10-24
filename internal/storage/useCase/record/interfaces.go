package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	// BulkCreateUpdate contracts bulk create/update operations.
	BulkCreateUpdate(context.Context, []*entity.Record) error
	// HeadList return all records where record head more than given head.
	HeadList(ctx context.Context, brnID entity.BranchID, head uint64) ([]*entity.Record, error)
}

type UseCase interface {
	// BulkCreateUpdate contracts bulk create/update operations.
	BulkCreateUpdate(context.Context, []*entity.Record) error
	// HeadList return all records where record head more than given head.
	HeadList(ctx context.Context, brnID entity.BranchID, head uint64) ([]*entity.Record, error)
}
