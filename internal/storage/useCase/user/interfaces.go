package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Create interface {
	Create(context.Context, *entity.User) error
}

//type Read interface {
//	Read(context.Context, locDomain.LocationInd) (*locDomain.Location, error)
//}
//
//type Update interface {
//	Update(context.Context, *locDomain.Location) error
//}
//
//type Delete interface {
//	Delete(context.Context, locDomain.LocationInd) error
//}

// CRUD контракт на реализацию операций CRUD: create, read, update, delete.
type CRUD interface {
	Create
	//Read
	//Update
	//Delete
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Repo interface {
	CRUD
	List
}

type List interface {
	List(ctx context.Context, offset, limit int) (*entity.UsersList, error)
}

type Add interface {
	Add(context.Context, *entity.User) error
}

//type Get interface {
//	Get(context.Context, locDomain.LocationInd) (*locDomain.Location, error)
//}
//
//type Set interface {
//	Set(context.Context, *locDomain.Location) error
//}
//
//type Del interface {
//	Del(context.Context, locDomain.LocationInd) error
//}

type UseCase interface {
	Add
	List
	//Get
	//Set
	//Del
}
