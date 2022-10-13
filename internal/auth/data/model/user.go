package model

import (
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/entity"
)

type User struct {
	gorm.Model
	Login string `gorm:"unique;not null"`
}

// DomainBind binds model type fields to entity model.
func (u *User) DomainBind() *entity.User {
	return &entity.User{
		ID:    entity.UserID(u.ID),
		Login: u.Login,
	}
}
