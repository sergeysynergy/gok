package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	// Create new record, pair: token - userID.
	Create(context.Context, *entity.Session) error
}

type UseCase interface {
	// Add new session record for given userID, return token value.
	Add(context.Context, entity.UserID) (*string, error)
}
