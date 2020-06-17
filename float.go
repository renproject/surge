package surge

import (
	"encoding/binary"
	"math"
)

const (
	// SizeHintF32 is the number of bytes required to represent a float32 value
	// in binary.
	SizeHintF32 = 4
	// SizeHintF64 is the number of bytes required to represent a float64 value
	// in binary.
	SizeHintF64 = 8
)

// MarshalF32 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalF32(x float32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintF32 || rem < SizeHintF32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint32(buf, math.Float32bits(x))
	return buf[SizeHintF32:], rem - SizeHintF32, nil
}

// MarshalF64 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalF64(x float64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintF64 || rem < SizeHintF64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint64(buf, math.Float64bits(x))
	return buf[SizeHintF64:], rem - SizeHintF64, nil
}

// UnmarshalF32 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalF32(x *float32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintF32 || rem < SizeHintF32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = math.Float32frombits(binary.BigEndian.Uint32(buf))
	return buf[SizeHintF32:], rem - SizeHintF32, nil
}

// UnmarshalF64 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalF64(x *float64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintF64 || rem < SizeHintF64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = math.Float64frombits(binary.BigEndian.Uint64(buf))
	return buf[SizeHintF64:], rem - SizeHintF64, nil
}
