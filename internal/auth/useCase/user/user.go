package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForUser describes all business-logic needed to work with `user` entity.
type UseCaseForUser struct {
	lg   *zap.Logger
	repo Repo
}

var _ UseCase = new(UseCaseForUser)

func New(logger *zap.Logger, repo Repo) *UseCaseForUser {
	s := &UseCaseForUser{
		lg:   logger,
		repo: repo,
	}

	return s
}

func (u *UseCaseForUser) SignIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
	errPrefix := "UseCaseForUser.SignIn"

	// Generate personal user encryption key.
	usr.KeyGen()

	id, err := u.repo.Create(ctx, usr)
	if err != nil {
		err = fmt.Errorf("%s - %w", errPrefix, err)
		u.lg.Error(err.Error())
		return nil, err
	}
	u.lg.Debug(fmt.Sprintf("New user with ID %d has been successfully created", id))

	sgdUser := &entity.SignedUser{
		Token: "bestTokenEver",
		Key:   usr.Key,
	}

	return sgdUser, nil
}

func (u *UseCaseForUser) SignIn(ctx context.Context, usr *entity.User) (*entity.SignedUser, error) {
