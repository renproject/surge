package surge

import (
	"reflect"
)

func sizeHintReflectedStruct(v reflect.Value) int {
	sizeHint := 0
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		if f := v.Field(i); f.IsValid() {
			sizeHint += sizeHintReflected(f)
		}
	}
	return sizeHint
}

func marshalReflectedStruct(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	var err error
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		if f := v.Field(i); f.IsValid() {
			if buf, rem, err = marshalReflected(f, buf, rem); err != nil {
				return buf, rem, err
			}
		}
	}
	return buf, rem, nil
}

func unmarshalReflectedStruct(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	var err error
	elem := v.Elem()
	numField := elem.NumField()
	for i := 0; i < numField; i++ {
		if f := elem.Field(i); f.IsValid() {
			if buf, rem, err = unmarshalReflected(f.Addr(), buf, rem); err != nil {
				return buf, rem, err
			}
		}
	}
	return buf, rem, nil
}
