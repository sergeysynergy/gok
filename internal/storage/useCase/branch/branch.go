package user

import (
	"context"
	"fmt"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
	recUC "github.com/sergeysynergy/gok/internal/storage/useCase/record"
	"go.uber.org/zap"
	"sync"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForBranch describes all business-logic needed to work with `branch` entity:
// - add new branch;
// - get branch data;
// - set branch data;
// - push branch and connected records;
// - pull branch and connected records.
type UseCaseForBranch struct {
	pushPullMu sync.RWMutex
	lg         *zap.Logger
	repo       Repo
	client     Client
	record     recUC.UseCase
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

	brn, err = u.repo.CreateReadByName(ctx, brn)
	if err != nil {
		return nil, err
	}
	u.lg.Debug(fmt.Sprintf("got branch: ID %d; name %s; head %d", brn.ID, brn.Name, brn.Head))

	return brn, nil
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

func (u *UseCaseForBranch) Push(ctx context.Context, token string, localBrn *entity.Branch, recs []*entity.Record) (*entity.Branch, error) {
	var err error
	logPrefix := "UseCaseForBranch.Push"
	u.lg.Debug(logPrefix)
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		err = gokErrors.ErrUserNotFound
		return nil, err
	}

	u.pushPullMu.Lock()
	defer u.pushPullMu.Unlock()

	// Get the latest branch head.
	localBrn.UserID = usr.ID
	freshBrn, err := u.repo.Read(ctx, localBrn)
	if err != nil {
		return nil, err
	}

	// Head check:
	if freshBrn.Head > localBrn.Head {
		return nil, gokErrors.ErrLocalBranchBehind
	}

	err = u.record.BulkCreateUpdate(ctx, recs)
	if err != nil {
		return nil, err
	}

	// IMPORTANT: push was successful - increase server branch head
	freshBrn.Head = localBrn.Head + 1
	err = u.repo.Update(ctx, freshBrn)
	if err != nil {
		return nil, err
	}

	u.lg.Debug(fmt.Sprintf("%s successful for branch: ID %d; name %s; head %d", logPrefix, freshBrn.ID, freshBrn.Name, freshBrn.Head))
	return freshBrn, nil
}

func (u *UseCaseForBranch) Pull(ctx context.Context, token string, localBrn *entity.Branch) (*entity.Branch, []*entity.Record, error) {
	u.lg.Debug("doing UseCaseForBranch.Pull")
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForBranch.Pull"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	usr, err := u.client.GetUser(ctx, token)
	if err != nil {
		err = gokErrors.ErrUserNotFound
		return nil, nil, err
	}

	u.pushPullMu.RLock()
	defer u.pushPullMu.RUnlock()

	// Get the latest branch head.
	localBrn.UserID = usr.ID
	freshBrn, err := u.repo.Read(ctx, localBrn)
	if err != nil {
		return nil, nil, err
	}

	if freshBrn.Head <= localBrn.Head {
		return nil, nil, gokErrors.ErrPullUpToDate
	}

	recs, err := u.record.HeadList(ctx, localBrn.ID, localBrn.Head)
	if err != nil {
		return nil, nil, err
	}

	u.lg.Debug(fmt.Sprintf("UseCaseForBranch.Pull - successfully pull %d records: branch `%s`; fresh head %d", len(recs), freshBrn.Name, freshBrn.Head))
	return freshBrn, recs, nil
}
