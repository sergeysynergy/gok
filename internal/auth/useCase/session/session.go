package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForSession is realization for session.UseCase interface.
// Describes all business-logic needed to work with `session` entity.
type UseCaseForSession struct {
	lg   *zap.Logger
	repo Repo
}

var _ UseCase = new(UseCaseForSession)

func New(logger *zap.Logger, repo Repo) *UseCaseForSession {
	s := &UseCaseForSession{
		lg:   logger,
		repo: repo,
	}

	return s
}

func (u *UseCaseForSession) Add(ctx context.Context, ses *entity.Session) error {
	errPrefix := "UseCaseForSession.Add"

	err := u.repo.Create(ctx, ses)
	if err != nil {
		err = fmt.Errorf("%s - %w", errPrefix, err)
		u.lg.Error(err.Error())
		return err
	}
	u.lg.Debug("New user session has been successfully created")

	return nil
}
