package model

import (
	"github.com/sergeysynergy/gok/internal/entity"
)

// Card is extension for basic record type to safely store credit card essentials.
type Card struct {
	ID      string `gorm:"primaryKey;not null"`
	Number  uint64 `gorm:"not null"`
	Code    uint64 `gorm:"not null"`
	Expired string `gorm:"not null"`
	Owner   string `gorm:"not null"`
}

func (m *Card) DomainBind() *entity.Card {
	return &entity.Card{
		ID:      entity.RecordID(m.ID),
		Number:  entity.NumberField(m.Number),
		Code:    entity.NumberField(m.Code),
		Expired: entity.StringField(m.Expired),
		Owner:   entity.StringField(m.Owner),
	}
}
