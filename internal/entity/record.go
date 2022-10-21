package entity

import (
	"fmt"
	"github.com/google/uuid"
	"time"

	gokConsts "github.com/sergeysynergy/gok/internal/consts"
)

type StringField interface {
	Encrypt(key string) error
	Decrypt(key string) (*string, error)
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

func NewRecord(
	key string,
	head uint64,
	branch string,
	description string,
	updatedAt time.Time,
	addition interface{},
) *Record {
	r := &Record{
		Head:        head,
		Branch:      branch,
		Description: Description(description),
		UpdatedAt:   updatedAt,
	}
	r.genID()

	err := r.Description.Encrypt(key)
	if err != nil {
		return nil
	}

	switch addition.(type) {
	default:
		r.Type = gokConsts.DESC
	}

	return r
}

func (r *Record) String() string {
	return fmt.Sprintf("%s\t %s\t %d\t %s\t %s", r.ID, r.Type, r.Head, r.UpdatedAt, r.Description)
}

// genID generate new record ID.
func (r *Record) genID() {
	r.ID = RecordID(uuid.New().String())
}

type Description string

// Make sure that the Description field implements the Field interface:
var _ StringField = new(Description)

func (d *Description) Encrypt(key string) error {
	return nil
}

func (d *Description) Decrypt(key string) (*string, error) {
	res := string(*d)
	return &res, nil
}
