// Package consts contains constants definition for all internal packages.
package consts

import "time"

const (
	// ServerGraceTimeout is time limit for gracefully shutdown service.
	ServerGraceTimeout = 10 * time.Second
)
