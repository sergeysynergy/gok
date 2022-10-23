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

func (r *Repo) Create(ctx context.Context, rec *entity.Record) (err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		recDB := model.Record{
			ID:          string(rec.ID),
			Head:        rec.Head,
			BranchID:    uint64(rec.BranchID),
			Description: string(rec.Description),
			UpdatedAt:   rec.UpdatedAt,
			Type:        string(rec.Type),
		}
		result := tx.Create(&recDB)
		err = result.Error
		if err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			return gokErrors.ErrRecordUnknown
		}

		switch ex := rec.Extension.(type) {
		case *entity.Text:
			textDB := model.Text{
				ID:   string(rec.ID), // always using base record ID
				Text: string(ex.Text),
			}
			err = tx.Create(&textDB).Error
			if err != nil {
				return err
			}
		default:
			if ex != nil {
				return gokErrors.ErrRecordUnknownExtensionType
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) Read(ctx context.Context, id entity.RecordID) (*entity.Record, error) {
	return nil, nil
}

func (r *Repo) Update(ctx context.Context, rec *entity.Record) (err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		recDB := model.Record{
			ID:          string(rec.ID),
			Head:        rec.Head,
			BranchID:    uint64(rec.BranchID),
			Description: string(rec.Description),
			UpdatedAt:   rec.UpdatedAt,
			Type:        string(rec.Type),
		}

		result := tx.Model(&recDB).Updates(&recDB)
		err = result.Error
		if err != nil {
			return err
		}
		if result.RowsAffected == 0 {
			err = gokErrors.ErrRecordNotFound
			return err
		}

		switch ex := rec.Extension.(type) {
		case *entity.Text:
			textDB := model.Text{
				ID:   string(rec.ID), // always using base record ID
				Text: string(ex.Text),
			}
			err = tx.Model(&textDB).Updates(&textDB).Error
			if err != nil {
				return err
			}
		default:
			if ex != nil {
				return gokErrors.ErrRecordUnknownExtensionType
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) addExtension(ctx context.Context, rec *entity.Record) (err error) {
	if rec.Type == gokConsts.DESC {
		return nil
	}

	tx := r.db.WithContext(ctx)
	switch rec.Type {
	case gokConsts.TEXT:
		textDB := &model.Text{ID: string(rec.ID)} // always using base record ID
		err = tx.Take(&textDB).Error
		if err != nil {
			return err
		}
		rec.Extension = textDB.DomainBind()
	default:
		return gokErrors.ErrRecordUnknownExtensionType
	}

	return nil
}

func (r *Repo) TypeList(ctx context.Context, brnID entity.BranchID, recType gokConsts.RecordType) ([]*entity.Record, error) {
	tx := r.db.WithContext(ctx)

	listDB := make([]*model.Record, 0)
	result := tx.Where("branch_id = ? AND type = ?", brnID, recType).Find(&listDB)
	err := result.Error
	if err != nil {
		return nil, err
	}

	list := make([]*entity.Record, 0, len(listDB))
	for _, v := range listDB {
		rec := v.DomainBind()

		// Search for extension for non DESC types
		err = r.addExtension(ctx, rec)
		if err != nil {
			return nil, err
		}

		list = append(list, rec)
	}

	return list, nil
}

// HeadList return all records where record head more than given head.
func (r *Repo) HeadList(ctx context.Context, brnID entity.BranchID, head uint64) ([]*entity.Record, error) {
	tx := r.db.WithContext(ctx)

	listDB := make([]*model.Record, 0)
	result := tx.Where("branch_id = ? AND head > ?", brnID, head).Find(&listDB)
	err := result.Error
	if err != nil && err != gorm.ErrRecordNotFound {
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

func (r *Repo) BulkCreateUpdate(ctx context.Context, recs []*entity.Record) (err error) {
	err = r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range recs {
			recDB := model.Record{
				ID:          string(v.ID),
				Head:        v.Head,
				BranchID:    uint64(v.BranchID),
				Description: string(v.Description),
				Type:        string(v.Type),
				UpdatedAt:   v.UpdatedAt,
			}

			result := tx.Take(&model.Record{ID: string(v.ID)})
			err = result.Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
			if result.RowsAffected > 0 {
				err = tx.Model(&recDB).Updates(&recDB).Error
				if err != nil {
					return err
				}
				continue
			}

			err = tx.Create(&recDB).Error
			if err != nil && err != gorm.ErrNotImplemented {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// ByIDsList return all records from given IDs slice.
func (r *Repo) ByIDsList(ctx context.Context, ids []string) ([]*entity.Record, error) {
	tx := r.db.WithContext(ctx)

	listDB := make([]*model.Record, 0, len(ids))
	result := tx.Where("id IN ?", ids).Find(&listDB)
	err := result.Error
	if err != nil && err != gorm.ErrRecordNotFound {
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
