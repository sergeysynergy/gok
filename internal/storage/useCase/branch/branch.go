package user

import (
	"context"
	"fmt"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForBranch describes all business-logic needed to work with `branch` entity:
// - add new branch;
// - get branch data.
type UseCaseForBranch struct {
	lg     *zap.Logger
	repo   Repo
	client Client
}

var _ UseCase = new(UseCaseForBranch)

func New(logger *zap.Logger, repo Repo, client Client) *UseCaseForBranch {
	s := &UseCaseForBranch{
		lg:     logger,
		repo:   repo,
		client: client,
	}

	return s
}

func (u *UseCaseForBranch) AddGet(ctx context.Context, token string) (*entity.Branch, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.Add"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	// TODO: move auth to interceptor
	usr, err := u.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		err = gokErrors.ErrUserNotFound
		return nil, err
	}

	brn := &entity.Branch{
		UserID: usr.ID,
		Name:   "default", // Limit using one branch so far.
	}

	brn, err = u.repo.CreateRead(ctx, brn)
	if err != nil {
		return nil, err
	}
	u.lg.Debug(fmt.Sprintf("Got branch: ID %d; name %s; head %d", brn.ID, brn.Name, brn.Head))

	return brn, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
/// Cross-service communication

// GetUser retrieve user data from `Auth` server by token.
func (u *UseCaseForBranch) GetUser(ctx context.Context, token string) (*entity.User, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.GetUser"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
