package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	// BulkCreateUpdate contracts bulk create/update operations.
	BulkCreateUpdate(context.Context, []*entity.Record) error
}

type UseCase interface {
	// BulkCreateUpdate contracts bulk create/update operations.
	BulkCreateUpdate(context.Context, []*entity.Record) error
}
