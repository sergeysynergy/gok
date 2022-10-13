package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

// CRUD контракт на реализацию операций CRUD: create, read, update, delete.
type CRUD interface {
	Create(context.Context, *entity.Session) error
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Repo interface {
	CRUD
}

type UseCase interface {
	Add(context.Context, *entity.Session) error
}
