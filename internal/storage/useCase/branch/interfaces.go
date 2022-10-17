package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

type Repo interface {
	// CreateRead create new branch record or read existing one.
	CreateRead(context.Context, *entity.Branch) (*entity.Branch, error)
}

type UseCase interface {
	// AddGet creates a new branch for user or return existing branch:
	// the user is identified through a token by cross-service request to `Auth` service.
	AddGet(ctx context.Context, token string) (*entity.Branch, error)
}

// Client defines contract for cross-service communication.
type Client interface {
	GetUser(ctx context.Context, token string) (*entity.User, error)
}
