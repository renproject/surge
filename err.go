package surge

import "fmt"

type ErrUnsupportedMarshalType struct {
	error
}

func newErrUnsupportedMarshalType(v interface{}) error {
	return ErrUnsupportedMarshalType{error: fmt.Errorf("marshal error: unsupported type %T", v)}
}

type ErrUnsupportedUnmarshalType struct {
	error
}

func newErrUnsupportedUnmarshalType(v interface{}) error {
	return ErrUnsupportedUnmarshalType{error: fmt.Errorf("unmarshal error: unsupported type %T", v)}
}

type ErrTagNotFound struct {
	error
}

func newErrTagNotFound(v uint16) error {
	return ErrTagNotFound{error: fmt.Errorf("unmarshal error: tag `surge:\"%v\"` not found", v)}
}

type ErrUnsupportedSizeHintType struct {
	error
}

func newErrUnsupportedSizeHintType(v interface{}) error {
	return ErrUnsupportedSizeHintType{error: fmt.Errorf("size hint error: unsupported type %T", v)}
}
