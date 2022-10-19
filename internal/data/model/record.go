// Package model contains types to use in database only.
package model

import (
	gokConsts "github.com/sergeysynergy/gok/internal/consts"
	"github.com/sergeysynergy/gok/internal/entity"
	"time"
)

// Record is basic type for all secret data to use in database.
type Record struct {
	ID          string    `gorm:"primaryKey;not null"`
	Head        uint64    `gorm:"not null"`
	Branch      string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (r *Record) DomainBind() *entity.Record {
	return &entity.Record{
		ID:          entity.RecordID(r.ID),
		Head:        r.Head,
		Branch:      r.Branch,
		Description: entity.Description(r.Description),
		Type:        gokConsts.RecordType(r.Type),
		UpdatedAt:   r.UpdatedAt,
	}
}
