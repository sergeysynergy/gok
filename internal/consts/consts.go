// Package consts contains constants definition for all internal packages.
package consts

import "time"

type RecordType string

const (
	// ServerGraceTimeout is time limit for gracefully shutdown service.
	ServerGraceTimeout = 10 * time.Second

	DESC = RecordType("DESC")
	TEXT = RecordType("TEXT")
	PASS = RecordType("PASS")
	CARD = RecordType("CARD")
	FILE = RecordType("FILE")
)
