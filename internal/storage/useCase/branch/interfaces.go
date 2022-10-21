package user

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
)

// Repo defines CRUD realization.
type Repo interface {
	Read(context.Context, *entity.Branch) (*entity.Branch, error)
	Update(context.Context, *entity.Branch) error
	// CreateRead create new branch record or read existing one.
	CreateRead(context.Context, *entity.Branch) (*entity.Branch, error)
}

type UseCase interface {
	// AddGet creates a new branch for user or return existing branch:
	// the user is identified through a token by cross-service request to `Auth` service.
	AddGet(ctx context.Context, token string) (*entity.Branch, error)
	// Get user branch info by user ID and branch name.
	Get(ctx context.Context, token string, brn *entity.Branch) (*entity.Branch, error)
	// Set updates user branch.
	Set(ctx context.Context, token string, brn *entity.Branch) error
	// Push updates form local branch to server.
	Push(ctx context.Context, token string, brn *entity.Branch, records []*entity.Record) (*entity.Branch, error)
	// Pull get records with more than local head for local update process.
	Pull(ctx context.Context, token string, brn *entity.Branch) (*entity.Branch, []*entity.Record, error)
}

// Client defines contract for cross-service communication.
type Client interface {
	GetUser(ctx context.Context, token string) (*entity.User, error)
}
