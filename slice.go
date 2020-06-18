package surge

import (
	"reflect"
)

func SizeHintBytes(v []byte) int {
	return 2 + len(v)
}

func MarshalBytes(v []byte, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintBytes(v) || rem < SizeHintBytes(v) {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	buf, rem, err := MarshalU16(uint16(len(v)), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	copy(buf, v)
	buf = buf[len(v):]
	rem -= len(v)
	return buf, rem, nil
}

func UnmarshalBytes(v *[]byte, buf []byte, rem int) ([]byte, int, error) {
	vLen := uint16(0)
	buf, rem, err := UnmarshalU16(&vLen, buf, rem)
	if err != nil {
		return buf, rem, err
	}

	if len(buf) < int(vLen) || rem < int(vLen) {
		return buf, rem, err
	}
	*v = make([]byte, vLen)
	copy(*v, buf)
	buf = buf[vLen:]
	rem -= int(vLen)

	return buf, rem, nil
}

func sizeHintReflectedSlice(v reflect.Value) int {
	sizeHint := 2
	for i := 0; i < v.Len(); i++ {
		sizeHint += sizeHintReflected(v.Index(i))
	}
	return sizeHint
}

func marshalReflectedSlice(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalU16(uint16(v.Len()), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	for i := 0; i < v.Len(); i++ {
		if buf, rem, err = marshalReflected(v.Index(i), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}

func unmarshalReflectedSlice(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	sliceLen := uint16(0)
	buf, rem, err := UnmarshalU16(&sliceLen, buf, rem)
	if err != nil {
		return buf, rem, err
	}

	elem := v.Elem()
	n := int(sliceLen)
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	n *= int(elem.Type().Elem().Size())
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	if rem < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	rem -= n

	elem.Set(reflect.MakeSlice(elem.Type(), int(sliceLen), int(sliceLen)))
	for i := uint16(0); i < sliceLen; i++ {
		if buf, rem, err = unmarshalReflected(elem.Index(int(i)).Addr(), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}
