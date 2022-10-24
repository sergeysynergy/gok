package useCase

import (
	"fmt"
	"time"

	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
)

// RecordAdd creates new record of given type.
func (u *GokUseCase) RecordAdd(cfg *entity.CLIConf, rec *entity.Record) error {
	var err error
	logPrefix := "GokUseCase.RecordAdd"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", logPrefix))
		}
	}()

	newRec := entity.NewRecord(
		cfg.Key,         // provide key for encryption
		cfg.LocalHead+1, // increase head counter for new records
		entity.BranchID(cfg.BranchID),
		rec.Description,
		time.Now(),
		rec.Extension,
	)

	err = u.repo.Create(u.ctx, newRec)
	if err != nil {
		return err
	}

	return nil
}

// RecordSet update existing record.
func (u *GokUseCase) RecordSet(cfg *entity.CLIConf, rec *entity.Record) error {
	var err error
	logPrefix := "GokUseCase.RecordSet"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", logPrefix))
		}
	}()

	if rec.ID == "" {
		return fmt.Errorf("empty record ID given")
	}

	// Create new record to apply encryption processing.
	updatedRec := entity.NewRecord(
		cfg.Key,         // provide key for encryption
		cfg.LocalHead+1, // increase head counter for updated records
		entity.BranchID(cfg.BranchID),
		rec.Description,
		time.Now(),
		rec.Extension, // nil value means default DESC type
	)
	// IMPORTANT: replace new record ID with existing one.
	updatedRec.ID = rec.ID

	err = u.repo.Update(u.ctx, updatedRec)
	if err != nil {
		return err
	}

	return nil
}

// RecordList return list records of given type for current local brand ID.
func (u *GokUseCase) RecordList(cfg *entity.CLIConf, recType gokConsts.RecordType) ([]*entity.Record, error) {
	var err error
	logPrefix := "GokUseCase.RecordList"
	defer func() {
		if err != nil {
			err = fmt.Errorf("%s - %w", logPrefix, err)
			u.lg.Error(err.Error())
		} else {
			u.lg.Debug(fmt.Sprintf("%s done successfully", logPrefix))
		}
	}()

	list, err := u.repo.TypeList(u.ctx, entity.BranchID(cfg.BranchID), recType)
	if err != nil {
		return nil, err
	}

	return list, nil
}
