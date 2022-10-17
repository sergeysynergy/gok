package model

import (
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/entity"
)

type Branch struct {
	gorm.Model
	UserID uint32 `gorm:"unique;not null"`
	Name   string `gorm:"not null"`
	Head   uint64 `gorm:"not null"`
}

// DomainBind binds model type fields to entity model.
func (b *Branch) DomainBind() *entity.Branch {
	return &entity.Branch{
		ID:     entity.BranchID(b.ID),
		UserID: entity.UserID(b.UserID),
		Name:   b.Name,
		Head:   b.Head,
	}
}
