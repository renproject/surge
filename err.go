package surge

import (
	"errors"
	"fmt"
)

// ErrUnexpectedEndOfBuffer is used when reading/writing from/to a buffer that
// has less space than expected.
var ErrUnexpectedEndOfBuffer = errors.New("unexpected end of buffer")

// ErrLengthOverflow is returned when the length of an array or slice has
// overflowed.
var ErrLengthOverflow = errors.New("max bytes exceeded")

// ErrUnsupportedMarshalType is returned when the an unsupported type is
// encountered during marshaling.
type ErrUnsupportedMarshalType struct {
	error
}

// NewErrUnsupportedMarshalType constructs a new unsupported marshal type error
// for the given type.
func NewErrUnsupportedMarshalType(v interface{}) error {
	return ErrUnsupportedMarshalType{error: fmt.Errorf("marshal error: unsupported type %T", v)}
}

// ErrUnsupportedUnmarshalType is returned when the an unsupported type is
// encountered during unmarshaling.
type ErrUnsupportedUnmarshalType struct {
	error
}

// NewErrUnsupportedUnmarshalType constructs a new unsupported unmarshal type
// error for the given type.
func NewErrUnsupportedUnmarshalType(v interface{}) error {
	return ErrUnsupportedUnmarshalType{error: fmt.Errorf("unmarshal error: unsupported type %T", v)}
}
