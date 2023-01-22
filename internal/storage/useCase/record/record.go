package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"

	"github.com/sergeysynergy/gok/internal/entity"
)

// UseCaseForRecord describes all business-logic needed to work with `record` entity at server side.
type UseCaseForRecord struct {
	lg   *zap.Logger
	repo Repo
}

var _ UseCase = new(UseCaseForRecord)

func New(logger *zap.Logger, repo Repo) *UseCaseForRecord {
	uc := &UseCaseForRecord{
		lg:   logger,
		repo: repo,
	}

	return uc
}

func (u *UseCaseForRecord) BulkCreateUpdate(ctx context.Context, recs []*entity.Record) error {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForRecord.BulkCreateUpdate"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	err = u.repo.BulkCreateUpdate(ctx, recs)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCaseForRecord) HeadList(ctx context.Context, brnID entity.BranchID, head uint64) ([]*entity.Record, error) {
	var err error
	defer func() {
		if err != nil {
			errPrefix := "UseCaseForRecord.HeadList"
			err = fmt.Errorf("%s - %w", errPrefix, err)
			u.lg.Error(err.Error())
		}
	}()

	recs, err := u.repo.HeadList(ctx, brnID, head)
	if err != nil {
		return nil, err
	}

	return recs, nil
}
