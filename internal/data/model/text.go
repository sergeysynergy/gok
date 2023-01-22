package model

import (
	"github.com/sergeysynergy/gok/internal/entity"
)

// Text is extension for Record model providing additional text field.
type Text struct {
	ID   string `gorm:"primaryKey;not null"`
	Text string `gorm:"not null"`
}

func (t *Text) DomainBind() *entity.Text {
	return &entity.Text{
		ID:   entity.RecordID(t.ID),
		Text: entity.StringField(t.Text),
	}
}
