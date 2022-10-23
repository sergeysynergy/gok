package entity

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	gokConsts "github.com/sergeysynergy/gok/internal/consts"
)

type RecordID string

// Record is basic type for all secret data.
type Record struct {
	ID          RecordID
	Head        uint64
	BranchID    BranchID
	Description StringField
	UpdatedAt   time.Time
	Type        gokConsts.RecordType

	Extension interface{}
}

func NewRecord(
	key string,
	head uint64,
	branchID BranchID,
	description StringField,
	updatedAt time.Time,
	//recType gokConsts.RecordType,
	extension interface{},
) *Record {
	r := &Record{
		Head:        head,
		BranchID:    branchID,
		Description: description,
		UpdatedAt:   updatedAt,
		//Type:        recType,
	}
	r.genID()

	err := r.Description.Encrypt(key)
	if err != nil {
		return nil
	}

	// Using new type functions for encryption processing.
	switch ex := extension.(type) {
	case *Text:
		r.Type = gokConsts.TEXT
		r.Extension = NewText(key, r.ID, ex.Text)
	default:
		// For default description type - extension should be nil.
		if extension != nil {
			log.Fatalln("Unknown record type")
		}
		r.Type = gokConsts.DESC
	}

	return r
}

func (r *Record) String() string {
	msg := fmt.Sprintf("%s\t %s\t %d\t %s", r.ID, r.Type, r.Head, r.Description)

	switch ex := r.Extension.(type) {
	case *Text:
		if ex != nil {
			msg = fmt.Sprintf("%s\t %s", msg, ex.Text)
		}
	}

	return msg
}

// genID generate new record ID.
func (r *Record) genID() {
	r.ID = RecordID(uuid.New().String())
}

type StringField string

func (d *StringField) Encrypt(key string) error {
	return nil
}

func (d *StringField) Decrypt(key string) (*string, error) {
	res := string(*d)
	return &res, nil
}
