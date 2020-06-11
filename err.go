package surge

import (
	"errors"
	"fmt"
)

// ErrUnexpectedEndOfBuffer is used when reading/writing from/to a buffer that
// has less spaced than expected.
var ErrUnexpectedEndOfBuffer = errors.New("unexpected end of buffer")

// ErrNotComparable is used when sorting maps that do not have comparable key
// types.
var ErrNotComparable = errors.New("not comparable")

// ErrMaxBytesExceeded is returned when too many bytes need to be allocated in
// memory.
var ErrMaxBytesExceeded = errors.New("max bytes exceeded")

// ErrLengthOverflow is returned when the length of an array or slice has
// overflowed.
var ErrLengthOverflow = errors.New("max bytes exceeded")

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
