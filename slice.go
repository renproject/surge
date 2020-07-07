package surge

import (
	"reflect"
)

// SizeHintBytes is the number of bytes required to represent the given byte
// slice in binary.
func SizeHintBytes(v []byte) int {
	return SizeHintU32 + len(v)
}

// MarshalBytes into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalBytes(v []byte, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalLen(uint32(len(v)), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	if len(buf) < len(v) || rem < len(v) {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	copy(buf, v)
	buf = buf[len(v):]
	rem -= len(v)
	return buf, rem, nil
}

// UnmarshalBytes from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalBytes(v *[]byte, buf []byte, rem int) ([]byte, int, error) {
	vLen := uint32(0)
	buf, rem, err := UnmarshalLen(&vLen, 1, buf, rem)
	if err != nil {
		return buf, rem, err
	}

	if len(buf) < int(vLen) {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*v = make([]byte, vLen)
	copy(*v, buf)
	buf = buf[vLen:]
	rem -= int(vLen)

	return buf, rem, nil
}

func sizeHintReflectedSlice(v reflect.Value) int {
	sizeHint := SizeHintU32
	for i := 0; i < v.Len(); i++ {
		sizeHint += sizeHintReflected(v.Index(i))
	}
	return sizeHint
}

func marshalReflectedSlice(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalLen(uint32(v.Len()), buf, rem)
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
	sliceLen := uint32(0)
	elem := v.Elem()
	size := int(elem.Type().Elem().Size())
	buf, rem, err := UnmarshalLen(&sliceLen, size, buf, rem)
	if err != nil {
		return buf, rem, err
	}
	rem -= int(sliceLen) * size

	elem.Set(reflect.MakeSlice(elem.Type(), int(sliceLen), int(sliceLen)))
	for i := uint32(0); i < sliceLen; i++ {
		if buf, rem, err = unmarshalReflected(elem.Index(int(i)).Addr(), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}
