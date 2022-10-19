package user

import (
	"context"
	"fmt"
	sesUC "github.com/sergeysynergy/gok/internal/auth/useCase/session"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForUser describes all business-logic needed to work with `user` entity.
type UseCaseForUser struct {
	lg      *zap.Logger
	repo    Repo
	session sesUC.UseCase
}

var _ UseCase = new(UseCaseForUser)

func New(logger *zap.Logger, repo Repo, session sesUC.UseCase) *UseCaseForUser {
	s := &UseCaseForUser{
		lg:      logger,
		repo:    repo,
		session: session,
	}

	return s
}

func (u *UseCaseForUser) SignIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForUser.SignIn"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	id, err := u.repo.Create(ctx, usr)
	if err != nil {
		return nil, err
	}
	u.lg.Debug(fmt.Sprintf("New user with ID %d has been successfully created", id))

	token, err := u.session.Add(ctx, id)
	if err != nil {
		return nil, err
	}

	sgdUser := entity.NewSignedUser(*token)

	return sgdUser, nil
}

func (u *UseCaseForUser) Login(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForUser.Login"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err = u.repo.Find(ctx, usr.Login)
	if err != nil {
		return nil, err
	}

	token, err := u.session.Add(ctx, usr.ID)
	if err != nil {
		return nil, err
	}

	sgdUser := entity.NewSignedUser(*token)

	u.lg.Debug(fmt.Sprintf("new login session for user ID %d has been created", usr.ID))
	return sgdUser, nil
}

func (u *UseCaseForUser) Get(ctx context.Context, token string) (*entity.User, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForUser.GetUser"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	// Get session by token, to reveal user ID.
	ses, err := u.session.Get(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user info.
	usr, err := u.repo.Read(ctx, ses.UserID)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
