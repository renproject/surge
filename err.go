package surge

import (
	"errors"
	"fmt"
)

// ErrMaxBytesExceeded is returned when too many bytes need to be allocated in
// memory.
var ErrMaxBytesExceeded = errors.New("max bytes exceeded")

// ErrUnsupportedMarshalType is returned when the an unsupported type is
// encountered during marshaling.
type ErrUnsupportedMarshalType struct {
	error
}

func newErrUnsupportedMarshalType(v interface{}) error {
	return ErrUnsupportedMarshalType{error: fmt.Errorf("marshal error: unsupported type %T", v)}
}

// ErrUnsupportedUnmarshalType is returned when the an unsupported type is
// encountered during unmarshaling.
type ErrUnsupportedUnmarshalType struct {
	error
}

func newErrUnsupportedUnmarshalType(v interface{}) error {
	return ErrUnsupportedUnmarshalType{error: fmt.Errorf("unmarshal error: unsupported type %T", v)}
}

// ErrBadLength is returned when unmarshaling into an array with the wrong
// length.
type ErrBadLength struct {
	error
}

func newErrBadLength(expected, got uint32) error {
	return ErrBadLength{error: fmt.Errorf("unmarshal error: expected len=%v, got len=%v", expected, got)}
}

// ErrNegativeLength is returned when unmarshaling into an array with the wrong
// length.
type ErrNegativeLength struct {
	error
}

func newErrNegativeLength(got int) error {
	return ErrNegativeLength{error: fmt.Errorf("unmarshal error: len>=0, got len=%v", got)}
}
