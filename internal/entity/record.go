package entity

import (
	"github.com/google/uuid"
	"time"

	gokConsts "github.com/sergeysynergy/gok/internal/consts"
)

type Field interface {
	Encrypt(key string) error
	Decrypt(key string) error
}

type RecordID string

// Record is basic type for all secret data.
type Record struct {
	ID          RecordID
	Head        uint64
	Branch      string
	Description Description
	Type        gokConsts.RecordType
	UpdatedAt   time.Time
}

func NewRecord(key string, head uint64, branch string, description string, updatedAt time.Time, addition interface{}) *Record {
	r := &Record{
		Head:        head,
		Branch:      branch,
		Description: Description(description),
		UpdatedAt:   updatedAt,
	}
	r.genID()
	r.Description.Encrypt(key)

	switch addition.(type) {
	default:
		r.Type = gokConsts.DESC
	}

	return r
}

// genID generate new record ID.
func (r *Record) genID() {
	r.ID = RecordID(uuid.New().String())
}

type Description string

// Make sure that the Description field implements the Field interface:
var _ Field = new(Description)

func (d *Description) Encrypt(key string) error {
	return nil
}

func (d *Description) Decrypt(key string) error {
	return nil
}
