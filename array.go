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
	arrayLen := v.Len()
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
	elem := v.Elem()
	arrayLen := elem.Len()
	if len(buf) < arrayLen || rem < arrayLen {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	var err error
	for i := 0; i < arrayLen; i++ {
		if buf, rem, err = unmarshalReflected(elem.Index(i).Addr(), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}
