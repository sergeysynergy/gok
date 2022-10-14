package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForSession describes all business-logic needed to work with `session` entity.
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

func (u *UseCaseForSession) Add(ctx context.Context, usrID entity.UserID) (*string, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForSession.Add"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	ses := entity.NewSession(usrID)

	err = u.repo.Create(ctx, ses)
	if err != nil {
		return nil, err
	}
	u.lg.Debug("New user session has been successfully created")

	return &ses.Token, nil
}
