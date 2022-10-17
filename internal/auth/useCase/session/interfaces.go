package session

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	// Create new session record, pair: token - userID.
	Create(context.Context, *entity.Session) error
	// Read session by given token.
	Read(ctx context.Context, token string) (*entity.Session, error)
}

type UseCase interface {
	// Add new session record for given userID, return token value.
	Add(context.Context, entity.UserID) (*string, error)
	// Get session by given token.
	Get(ctx context.Context, token string) (*entity.Session, error)
}
