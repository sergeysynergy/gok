package branch

import (
	"context"
	"github.com/sergeysynergy/gok/internal/storage/data/model"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/entity"
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

func (r *Repo) CreateRead(ctx context.Context, brn *entity.Branch) (branch *entity.Branch, err error) {
	tx := r.db.WithContext(ctx)

	brnDB := model.Branch{
		UserID: uint32(brn.UserID),
		Name:   brn.Name,
	}

	result := tx.Where("user_id = ? AND name = ?", brn.UserID, brn.Name).Take(&brnDB)
	err = result.Error
	if err != nil {
		result = tx.Create(&brnDB)
		err = result.Error
		if err != nil {
			return nil, err
		}
	}

	return brnDB.DomainBind(), nil
}

func (r *Repo) Read(ctx context.Context, brn *entity.Branch) (branch *entity.Branch, err error) {
	tx := r.db.WithContext(ctx)

	brnDB := model.Branch{
		UserID: uint32(brn.UserID),
		Name:   brn.Name,
	}

	result := tx.Where("user_id = ? AND name = ?", brn.UserID, brn.Name).Take(&brnDB)
	err = result.Error
	if err != nil {
		return nil, err
	}

	return brnDB.DomainBind(), nil
}

func (r *Repo) Update(ctx context.Context, brn *entity.Branch) error {
	tx := r.db.WithContext(ctx)

	brnDB := model.Branch{
		UserID:     uint32(brn.UserID),
		Name:       brn.Name,
		ServerHead: brn.Head,
	}

	result := tx.Where("user_id = ? AND name = ?", brn.UserID, brn.Name).Updates(&brnDB)
	err := result.Error
	if err != nil {
		return err
	}

	return nil
}
