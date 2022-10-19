package record

import (
	"context"
	"gorm.io/gorm"

	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/data/model"
	"github.com/sergeysynergy/gok/internal/entity"
	gokErrors "github.com/sergeysynergy/gok/internal/errors"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	r := &Repo{
		db: db,
	}
	return r
}

func (r *Repo) Create(ctx context.Context, rec *entity.Record) error {
	tx := r.db.WithContext(ctx)

	recDB := model.Record{
		ID:          string(rec.ID),
		Head:        rec.Head,
		Branch:      rec.Branch,
		Description: string(rec.Description),
		Type:        string(rec.Type),
		UpdatedAt:   rec.UpdatedAt,
	}

	result := tx.Create(&recDB)
	err := result.Error
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		err = gokErrors.ErrRecordUnknown
		return err
	}

	return nil
}

func (r *Repo) Read(ctx context.Context, id entity.RecordID) (*entity.Record, error) {
	//tx := r.db.WithContext(ctx)
	//
	//usrDB := model.User{}
	//err := tx.Take(&usrDB, id).Error
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, fmt.Errorf("%s: %w", err, gokErrors.ErrUserNotFound)
	//	}
	//	return nil, err
	//}
	//
	//return usrDB.DomainBind(), nil

	return nil, nil
}

func (r *Repo) Update(ctx context.Context, rec *entity.Record) error {
	tx := r.db.WithContext(ctx)

	recDB := model.Record{
		ID:          string(rec.ID),
		Head:        rec.Head,
		Branch:      rec.Branch,
		Description: string(rec.Description),
		Type:        string(rec.Type),
		UpdatedAt:   rec.UpdatedAt,
	}

	result := tx.Model(&recDB).Updates(&recDB)
	err := result.Error
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		err = gokErrors.ErrRecordNotFound
		return err
	}

	return nil
}

func (r *Repo) List(ctx context.Context, recType gokConsts.RecordType) ([]*entity.Record, error) {
	tx := r.db.WithContext(ctx)

	listDB := make([]*model.Record, 0)
	result := tx.Where("type = ?", recType).Find(&listDB)
	err := result.Error
	if err != nil {
		return nil, err
	}

	list := make([]*entity.Record, 0, len(listDB))
	for _, v := range listDB {
		list = append(list, v.DomainBind())
	}

	return list, nil
}

func (r *Repo) ListForPush(ctx context.Context, localHead uint64) ([]*entity.Record, error) {
	tx := r.db.WithContext(ctx)

	listDB := make([]*model.Record, 0)
	result := tx.Where("head > ?", localHead).Find(&listDB)
	err := result.Error
	if err != nil {
		return nil, err
	}
	if result.RowsAffected == 0 {
		err = gokErrors.ErrRecordNotFound
		return nil, err
	}

	list := make([]*entity.Record, 0, len(listDB))
	for _, v := range listDB {
		list = append(list, v.DomainBind())
	}

	return list, nil
}

func (r *Repo) BulkCreateUpdate(ctx context.Context, recs []*entity.Record) error {
	tx := r.db.WithContext(ctx)
	for _, v := range recs {
		recDB := model.Record{
			ID:          string(v.ID),
			Head:        v.Head,
			Branch:      v.Branch,
			Description: string(v.Description),
			Type:        string(v.Type),
			UpdatedAt:   v.UpdatedAt,
		}

		err := tx.Create(&recDB).Error
		if err != nil {
			err = tx.Model(&recDB).Updates(&recDB).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
