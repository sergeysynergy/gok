package useCase

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"

	gokClient "github.com/sergeysynergy/gok/internal/cli/client"
	"github.com/sergeysynergy/gok/internal/entity"
)

type GokUseCase struct {
	lg     *zap.Logger
	ctx    context.Context
	client Client
}

func New(logger *zap.Logger, client *gokClient.Client) *GokUseCase {
	const (
		defaultTimeout = 10 * time.Second
	)
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeout)

	uc := &GokUseCase{
		ctx:    ctx,
		lg:     logger,
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
