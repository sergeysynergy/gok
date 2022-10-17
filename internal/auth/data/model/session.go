package model

import "github.com/sergeysynergy/gok/internal/entity"

type Session struct {
	UserID int32  `gorm:"not null;uniqueIndex:UserIDAndToken"`
	Token  string `gorm:"not null;uniqueIndex:UserIDAndToken"`
}

// DomainBind binds model type fields to entity model.
func (s *Session) DomainBind() *entity.Session {
	return &entity.Session{
		UserID: entity.UserID(s.UserID),
		Token:  s.Token,
	}
}
