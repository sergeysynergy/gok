package session

import (
	"context"
	"fmt"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
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

func (u *UseCaseForSession) Get(ctx context.Context, token string) (*entity.Session, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForSession.Get"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	ses, err := u.repo.Read(ctx, token)
	if err != nil {
		return nil, err
	}
	// Catch for possible zero user id.
	if ses.UserID == 0 {
		err = gokErrors.ErrUserZeroID
		return nil, err
	}

	u.lg.Debug(fmt.Sprintf("Got session: userID %d", ses.UserID))

	return ses, nil
}
