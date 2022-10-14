package user

import (
	"context"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/auth/data/model"
	"github.com/sergeysynergy/gok/internal/entity"
	serviceErrors "github.com/sergeysynergy/gok/internal/errors"
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
		return 0, serviceErrors.ErrUserAlreadyExists
	}

	id = entity.UserID(usrDB.ID)

	return id, nil
}

//func (r *Repo) Read(ctx context.Context, id int32) (*entity.User, error) {
//	tx := r.db.WithContext(ctx)
//
//	usrDB := model.User{}
//	err := tx.Take(&usrDB, id).Error
//	if err != nil {
//		return nil, err
//	}
//
//	return &entity.User{}, nil
//}

//func (r *Repo) List(ctx context.Context, offset, limit int) (*entity.UsersList, error) {
//	tx := r.db.WithContext(ctx)
//
//	var usersDB []*model.User
//	err := tx.
//		Order("ID").
//		Offset(offset).
//		Limit(limit).
//		Find(&usersDB).
//		Error
//	if err != nil {
//		return nil, err
//	}
//
//	users := make([]*entity.User, 0, len(usersDB))
//	for _, v := range usersDB {
//		usr := v.DomainBind()
//		users = append(users, usr)
//	}
//
//	return &entity.UsersList{
//		Users: users,
//	}, nil
//}
