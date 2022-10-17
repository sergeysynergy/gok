package user

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/auth/data/model"
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

func (r *Repo) Read(ctx context.Context, token string) (*entity.Session, error) {
	tx := r.db.WithContext(ctx)
	sesDB := model.Session{}

	result := tx.Where("token = ?", token).Find(&sesDB)
	err := result.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", err, gokErrors.ErrSessionNotFound)
		}
		return nil, err
	}
	if result.RowsAffected == 0 {
		return nil, gokErrors.ErrSessionNotFound
	}

	return sesDB.DomainBind(), nil
}
