package useCase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	gokClient "github.com/sergeysynergy/gok/internal/cli/delivery/client"
	recRepo "github.com/sergeysynergy/gok/internal/data/repository/sql/record"
	"github.com/sergeysynergy/gok/internal/entity"
)

type GokUseCase struct {
	lg     *zap.Logger
	ctx    context.Context
	repo   Repo
	client Client
}

func New(logger *zap.Logger, repo *recRepo.Repo, client *gokClient.Client) *GokUseCase {
	const (
		defaultTimeout = 10 * time.Second
	)
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeout)

	uc := &GokUseCase{
		ctx:    ctx,
		lg:     logger,
		repo:   repo,
		client: client,
	}

	return uc
}

func (u *GokUseCase) SignIn(usrCLI *entity.CLIUser) (*entity.SignedUser, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.SignIn"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr := &entity.User{
		Login: usrCLI.Login,
	}

	signedUsr, err := u.client.SignIn(u.ctx, usr)
	if err != nil {
		return nil, err
	}

	return signedUsr, nil
}

func (u *GokUseCase) Login(usrCLI *entity.CLIUser) (*entity.SignedUser, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.Login"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr := &entity.User{
		Login: usrCLI.Login,
	}

	signedUsr, err := u.client.Login(u.ctx, usr)
	if err != nil {
		return nil, err
	}

	return signedUsr, nil
}

func (u *GokUseCase) Init(token string) (*entity.Branch, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.Init"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	brn, err := u.client.Init(u.ctx, token)
	if err != nil {
		return nil, err
	}

	return brn, nil
}

func (u *GokUseCase) Push(token string, branch string, head uint64) (*entity.Branch, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.Push"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	records, err := u.repo.ListForPush(u.ctx, head)
	if err != nil {
		return nil, err
	}

	brn := &entity.Branch{
		Name: branch,
		Head: head,
	}

	brn, err = u.client.Push(u.ctx, token, brn, records)
	if err != nil {
		return nil, err
	}

	return brn, nil
}
