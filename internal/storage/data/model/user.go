package model

import (
	"github.com/sergeysynergy/gok/internal/entity"
)

type User struct {
	ID    int32  `gorm:"primaryKey;column:ind;type:int(11);not null" json:"id"`
	Login string `gorm:"unique;column:name;type:varchar(255);not null" json:"login"`
}

// TableName get sql table name.
func (u *User) TableName() string {
	return "user"
}

// DomainBind binds model type fields to entity model.
func (u *User) DomainBind() *entity.User {
	return &entity.User{
		ID:    entity.UserID(u.ID),
		Login: u.Login,
	}
}

// UserColumns get sql column name.
var UserColumns = struct {
	ID    string
	Login string
}{
	ID:    "id",
	Login: "login",
}
