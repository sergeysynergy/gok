package user

import (
	"context"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/auth/data/model"
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

func (r *Repo) Create(ctx context.Context, ses *entity.Session) error {
	tx := r.db.WithContext(ctx)

	sesDB := model.Session{
		UserID: int32(ses.UserID),
		Token:  ses.Token,
	}

	err := tx.Create(&sesDB).Error
	if err != nil {
		return err
	}

	return nil
}
