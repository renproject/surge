package surge

import (
	"strings"
)

// SizeHintString is the number of bytes required to represent the given string
// in binary.
func SizeHintString(v string) int {
	return SizeHintU32 + len(v)
}

// MarshalString into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalString(v string, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalLen(uint32(len(v)), buf, rem)
	if err != nil {
		return buf, rem, err
	}
	if len(buf) < len(v) || rem < len(v) {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	n := copy(buf, v)
	return buf[n:], rem - n, nil
}

// UnmarshalString from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalString(v *string, buf []byte, rem int) ([]byte, int, error) {
	strLen := uint32(0)
	buf, rem, err := UnmarshalLen(&strLen, 1, buf, rem)
	if err != nil {
		return buf, rem, err
	}
	n := int(strLen)
	if len(buf) < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	b := strings.Builder{}
	if _, err := b.Write(buf[:n]); err != nil {
		return buf[:n], rem - n, err
	}
	*v = b.String()
	return buf[n:], rem - n, nil
}
