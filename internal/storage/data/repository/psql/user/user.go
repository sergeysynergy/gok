package user

import (
	"context"
	"fmt"
	usrModel "github.com/sergeysynergy/gok/internal/auth/data/model"
	"github.com/sergeysynergy/gok/internal/entity"
	serverErrors "github.com/sergeysynergy/gok/internal/errors"
	"github.com/sergeysynergy/gok/pkg/utils"
	"gorm.io/gorm"
)

type Repo struct {
	DB *gorm.DB
}

type Option func(repo *Repo)

func New(gormDB *gorm.DB, opts ...Option) *Repo {
	r := &Repo{
		DB: gormDB,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (r Repo) Create(ctx context.Context, usr *entity.User) (err error) {
	if usr == nil {
		return serverErrors.ErrUserInvalidArgument
	}

	usrDB := usrModel.User{
		ID:    int32(usr.ID),
		Login: usr.Login,
	}

	tx := r.DB.WithContext(ctx)
	err = tx.Transaction(func(txx *gorm.DB) error {
		err = tx.Create(&usrDB).Error
		if err != nil {
			if utils.IsMySQLDuplicateError(err) {
				return fmt.Errorf("%w: %s", serverErrors.ErrUserAlreadyExists, err)
			}
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) List(ctx context.Context, offset, limit int) (*entity.UsersList, error) {
	tx := r.DB.WithContext(ctx)

	var usersDB []*usrModel.User
	err := tx.
		Order(usrModel.UserColumns.ID).
		Offset(offset).
		Limit(limit).
		Find(&usersDB).
		Error
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, 0, len(usersDB))
	for _, v := range usersDB {
		usr := v.DomainBind()
		users = append(users, usr)
	}

	return &entity.UsersList{
		Users: users,
	}, nil
}
