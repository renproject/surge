package surge

import (
	"reflect"
)

func sizeHintReflectedArray(v reflect.Value) int {
	sizeHint := 0
	for i := 0; i < v.Len(); i++ {
		sizeHint += sizeHintReflected(v.Index(i))
	}
	return sizeHint
}

func marshalReflectedArray(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	arrayLen := v.Type().Len()
	if len(buf) < arrayLen || rem < arrayLen {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	var err error
	for i := 0; i < arrayLen; i++ {
		if buf, rem, err = marshalReflected(v.Index(i), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}

func unmarshalReflectedArray(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	arrayLen := v.Elem().Len()
	if len(buf) < arrayLen || rem < arrayLen {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	var err error
	for i := 0; i < arrayLen; i++ {
		if buf, rem, err = unmarshalReflected(v.Index(i).Addr(), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}
