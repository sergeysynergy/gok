package model

import (
	"github.com/sergeysynergy/gok/internal/entity"
)

// File is extension for basic record type to store binary data.
type File struct {
	ID   string `gorm:"primaryKey;not null"`
	File []byte `gorm:"not null"`
}

func (m *File) DomainBind() *entity.File {
	return &entity.File{
		ID:   entity.RecordID(m.ID),
		File: m.File,
	}
}
