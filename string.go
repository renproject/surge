package surge

import (
	"reflect"
	"strings"
)

func SizeHintString(v string) int {
	return 2 + len(v)
}

func MarshalString(v string, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalU16(uint16(len(v)), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	if len(buf) < len(v) || rem < len(v) {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	n := copy(buf, v)
	return buf[n:], rem - n, nil
}

func UnmarshalString(v *string, buf []byte, rem int) ([]byte, int, error) {
	strLen := uint16(0)
	buf, rem, err := UnmarshalU16(&strLen, buf, rem)
	if err != nil {
		return buf, rem, err
	}
	n := int(strLen)
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	if len(buf) < n || rem < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	b := strings.Builder{}
	if _, err := b.Write(buf[:n]); err != nil {
		return buf[:n], rem - n, err
	}
	*v = b.String()
	return buf[n:], rem - n, nil
}

func sizeHintReflectedString(v reflect.Value) int {
	return 2 + v.Len()
}

func marshalReflectedString(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalU16(uint16(v.Len()), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	if len(buf) < v.Len() || rem < v.Len() {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	n := copy(buf, v.String())
	return buf[n:], rem - n, nil
}

func unmarshalReflectedString(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	strLen := uint16(0)
	buf, rem, err := UnmarshalU16(&strLen, buf, rem)
	if err != nil {
		return buf, rem, err
	}
	n := int(strLen)
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	if len(buf) < n || rem < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	b := strings.Builder{}
	if _, err := b.Write(buf[:n]); err != nil {
		return buf[:n], rem - n, err
	}
	v.Elem().SetString(b.String())
	return buf[n:], rem - n, nil
}
