package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

// CRUD contract for realization of create, read, update and delete methods.
type CRUD interface {
	Create(context.Context, *entity.User) (entity.UserID, error)
}

// Repo contract methods to work with `user` repository.
type Repo interface {
	CRUD
}

// UseCase contract methods to work with `user` entity.
type UseCase interface {
	SignIn(context.Context, *entity.User) (*entity.SignedUser, error)
}
