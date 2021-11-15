package models

import (
	"fmt"
)

const (
	// Error Codes
	ErrFailedConversionCode = iota + 1000
	ErrProcessingRequestCode
)

// Error represents a deployer error.
type Error struct {
	Code    int
	Message string
}

// Error returns the string representation of a deployer's Error.
func (err Error) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}
