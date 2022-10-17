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

func (r *Repo) Create(ctx context.Context, usr *entity.User) (id entity.UserID, err error) {
	tx := r.db.WithContext(ctx)

	usrDB := model.User{
		Login: usr.Login,
	}

	result := tx.Create(&usrDB)
	if result.Error != nil {
		if err != nil {
			return 0, err
		}
	}
	if result.RowsAffected == 0 {
		return 0, gokErrors.ErrUserAlreadyExists
	}

	id = entity.UserID(usrDB.ID)

	return id, nil
}

func (r *Repo) Read(ctx context.Context, id entity.UserID) (*entity.User, error) {
	if id == 0 {
		return nil, gokErrors.ErrUserInvalidArgument
	}

	tx := r.db.WithContext(ctx)

	usrDB := model.User{}
	err := tx.Take(&usrDB, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%s: %w", err, gokErrors.ErrUserNotFound)
		}
		return nil, err
	}

	return usrDB.DomainBind(), nil
}
