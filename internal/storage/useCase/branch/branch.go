package user

import (
	"context"
	"fmt"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	recUC "github.com/sergeysynergy/gok/internal/storage/useCase/record"
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
	record recUC.UseCase
}

var _ UseCase = new(UseCaseForBranch)

func New(logger *zap.Logger, repo Repo, client Client, record recUC.UseCase) *UseCaseForBranch {
	s := &UseCaseForBranch{
		lg:     logger,
		repo:   repo,
		client: client,
		record: record,
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

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
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

func (u *UseCaseForBranch) Push(ctx context.Context, token string, localBrn *entity.Branch, records []*entity.Record) (*entity.Branch, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.Push"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		err = gokErrors.ErrUserNotFound
		return nil, err
	}

	localBrn.UserID = usr.ID
	freshBrn, err := u.repo.Read(ctx, localBrn)
	if err != nil {
		return nil, err
	}

	// Head check:
	if freshBrn.Head > localBrn.Head {
		return nil, gokErrors.ErrLocalBranchBehind
	}

	err = u.record.BulkCreateUpdate(ctx, records)
	if err != nil {
		return nil, err
	}

	// Push was successful: increase server branch head
	freshBrn.Head = localBrn.Head + 1
	err = u.repo.Update(ctx, freshBrn)
	if err != nil {
		return nil, err
	}

	u.lg.Debug(fmt.Sprintf("push successful: branch %s; server head %d", freshBrn.Name, freshBrn.Head))
	return freshBrn, nil
}

func (u *UseCaseForBranch) Get(ctx context.Context, token string, brn *entity.Branch) (*entity.Branch, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.Get"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		err = gokErrors.ErrUserNotFound
		return nil, err
	}
	brn.UserID = usr.ID

	freshBrn, err := u.repo.Read(ctx, brn)
	if err != nil {
		return nil, err
	}

	return freshBrn, nil
}

func (u *UseCaseForBranch) Set(ctx context.Context, token string, brn *entity.Branch) error {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.Get"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		err = gokErrors.ErrUserNotFound
		return err
	}
	brn.UserID = usr.ID

	err = u.repo.Update(ctx, brn)
	if err != nil {
		return err
	}

	return nil
}
