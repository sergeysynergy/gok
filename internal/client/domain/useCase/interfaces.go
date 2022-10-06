package useCase

import (
	"context"

	locDomain "back-mfsb/service/location/internal/domain/entity"
)

type Create interface {
	Create(context.Context, *locDomain.Location) (*locDomain.LocationInd, error)
}

type Read interface {
	Read(context.Context, locDomain.LocationInd) (*locDomain.Location, error)
}

type Update interface {
	Update(context.Context, *locDomain.Location) error
}

type Delete interface {
	Delete(context.Context, locDomain.LocationInd) error
}

// CRUD контракт на реализацию операций CRUD: create, read, update, delete.
type CRUD interface {
	Create
	Read
	Update
	Delete
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ByInd interface {
	ByInd(context.Context) (locDomain.ByInd, error)
}

type SortByInd interface {
	SortByInd(ctx context.Context, offset, limit int) (locDomain.SortByInd, error)
}

type Repo interface {
	CRUD
	ByInd
	SortByInd
}

type Add interface {
	Add(context.Context, *locDomain.Location) error
}

type Get interface {
	Get(context.Context, locDomain.LocationInd) (*locDomain.Location, error)
}

type Set interface {
	Set(context.Context, *locDomain.Location) error
}

type Del interface {
	Del(context.Context, locDomain.LocationInd) error
}

type UseCase interface {
	Add
	Get
	Set
	Del
	ByInd
}
