package model

import (
	"github.com/sergeysynergy/gok/internal/entity"
)

// Pass is extension for basic record type with login and password string types.
type Pass struct {
	ID   string `gorm:"primaryKey;not null"`
	Login string `gorm:"not null"`
	Password string `gorm:"not null"`
}

func (m *Pass) DomainBind() *entity.Pass {
	return &entity.Pass{
		ID:   entity.RecordID(m.ID),
		Login: entity.StringField(m.Login),
		Password: entity.StringField(m.Password),
	}
}
