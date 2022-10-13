package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForUser is realization for user.UseCase interface.
// Describes all business-logic needed to work with `user` entity.
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
func (u *UseCaseForUser) Add(ctx context.Context, loc *entity.User) (err error) {
	errPrefix := "UserUseCase.Add"

	err = u.repo.Create(ctx, loc)
	if err != nil {
		return fmt.Errorf("%s - %w", errPrefix, err)
	}

	u.lg.Debug("New user has been successfully created")

	return nil
}

func (u *UseCaseForUser) List(ctx context.Context, offset, limit int) (*entity.UsersList, error) {
	errPrefix := "UserUseCase.List"

	list, err := u.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s - %w", errPrefix, err)
	}

	return list, nil
}

//func (p *PersonUseCase) Del(ctx context.Context, psnInd psnDomain.PersonInd) error {
//errPrefix := "useCase Person.Del"

//err := p.repo.Delete(ctx, psnInd)
//if err != nil {
//	return fmt.Errorf("%s - %w", errPrefix, err)
//}

//	return nil
//}

//func (l *LocationUseCase) Get(ctx context.Context, locInd locDomain.LocationInd) (*locDomain.Location, error) {
//	errPrefix := "useCase Location.Get"
//
//	loc, err := l.repo.Read(ctx, locInd)
//	if err != nil {
//		return nil, fmt.Errorf("%s - %w", errPrefix, err)
//	}
//
//	return loc, nil
//}

//func (l *LocationUseCase) Set(ctx context.Context, loc *locDomain.Location) error {
//	errPrefix := "useCase Location.Set"
//
//	err := l.repo.Update(ctx, loc)
//	if err != nil {
//		return fmt.Errorf("%s - %w", errPrefix, err)
//	}
//
//	return nil
//}
