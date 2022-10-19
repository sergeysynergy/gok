package useCase

import (
	"fmt"
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

// DescAdd creates new record of DESC type: contains just description besides meta fields.
func (u *GokUseCase) DescAdd(rec *entity.Record) error {
	var err error
	defer func() {
		prefix := "GokUseCase.DescAdd"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	err = u.repo.Create(u.ctx, rec)
	if err != nil {
		return err
	}

	return nil
}

// DescSet update existing record of DESC type.
func (u *GokUseCase) DescSet(rec *entity.Record) error {
	var err error
	defer func() {
		prefix := "GokUseCase.DescSet"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	err = u.repo.Update(u.ctx, rec)
	if err != nil {
		return err
	}

	return nil
}

// DescList return list records of DESC type.
func (u *GokUseCase) DescList() ([]*entity.Record, error) {
	var err error
	defer func() {
		prefix := "GokUseCase.DescList"
		if err != nil {
			err = fmt.Errorf("%s - %w", prefix, err)
			u.lg.Error(err.Error())
		}
	}()

	list, err := u.repo.List(u.ctx, gokConsts.DESC)
	if err != nil {
		return nil, err
	}

	return list, nil
}
