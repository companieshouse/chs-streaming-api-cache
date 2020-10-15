// Package offset provides functions for validating request offsets.
package offset

import (
	"errors"
	"strconv"
)

var (
	// ErrNegativeOffset is the error returned when the requested offset is less than zero.
	ErrNegativeOffset = errors.New("offset should be greater than zero")
	// ErrOutOfRange is the error returned when the requested offset is out of range.
	ErrOutOfRange = errors.New("requested offset is out of range")
)

// Offset - structure for offset
type Offset struct {
}

// Interface - interface for offset methods
type Interface interface {
	Parse(offset string) (int64, error)
}

// NewOffset - return a new Offset
func NewOffset() Interface {
	return &Offset{}
}

// Parse parses a string offset as an int64.
func (of *Offset) Parse(offset string) (int64, error) {
	if offset == "" {
		return 0, nil
	}

	o, err := strconv.Atoi(offset)
	if err != nil {
		return 0, err
	}
	if o < 0 {
		return 0, ErrNegativeOffset
	}

	return int64(o), nil
}
