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

//func (u *UseCaseForUser) LogIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
//	return nil, nil
//}
