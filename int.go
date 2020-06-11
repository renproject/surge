package surge

import (
	"encoding/binary"
)

const (
	// SizeHintU8 is the number of bytes required to represent a uint8 value in
	// binary.
	SizeHintU8 = 1
	// SizeHintU16 is the number of bytes required to represent a uint16 value
	// in binary.
	SizeHintU16 = 2
	// SizeHintU32 is the number of bytes required to represent a uint32 value
	// in binary.
	SizeHintU32 = 4
	// SizeHintU64 is the number of bytes required to represent a uint64 value
	// in binary.
	SizeHintU64 = 8
	// SizeHintI8 is the number of bytes required to represent a int8 value in
	// binary.
	SizeHintI8 = 1
	// SizeHintI16 is the number of bytes required to represent a int16 value in
	// binary.
	SizeHintI16 = 2
	// SizeHintI32 is the number of bytes required to represent a int32 value in
	// binary.
	SizeHintI32 = 4
	// SizeHintI64 is the number of bytes required to represent a int64 value in
	// binary.
	SizeHintI64 = 8
)

// MarshalU8 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalU8(x uint8, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU8 || rem < SizeHintU8 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	buf[0] = x
	return buf[SizeHintU8:], rem - SizeHintU8, nil
}

// MarshalU16 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalU16(x uint16, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU16 || rem < SizeHintU16 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint16(buf, x)
	return buf[SizeHintU16:], rem - SizeHintU16, nil
}

// MarshalU32 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalU32(x uint32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU32 || rem < SizeHintU32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint32(buf, x)
	return buf[SizeHintU32:], rem - SizeHintU32, nil
}

// MarshalU64 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalU64(x uint64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU64 || rem < SizeHintU64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint64(buf, x)
	return buf[SizeHintU64:], rem - SizeHintU64, nil
}

// MarshalI8 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalI8(x int8, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI8 || rem < SizeHintI8 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	buf[0] = uint8(x)
	return buf[SizeHintI8:], rem - SizeHintI8, nil
}

// MarshalI16 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalI16(x int16, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI16 || rem < SizeHintI16 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint16(buf, uint16(x))
	return buf[SizeHintI16:], rem - SizeHintI16, nil
}

// MarshalI32 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalI32(x int32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI32 || rem < SizeHintI32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint32(buf, uint32(x))
	return buf[SizeHintI32:], rem - SizeHintI32, nil
}

// MarshalI64 into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalI64(x int64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI64 || rem < SizeHintI64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	binary.BigEndian.PutUint64(buf, uint64(x))
	return buf[SizeHintI64:], rem - SizeHintI64, nil
}

// UnmarshalU8 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalU8(x *uint8, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU8 || rem < SizeHintU8 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = buf[0]
	return buf[SizeHintU8:], rem - SizeHintU8, nil
}

// UnmarshalU16 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalU16(x *uint16, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU16 || rem < SizeHintU16 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = binary.BigEndian.Uint16(buf[:])
	return buf[SizeHintU16:], rem - SizeHintU16, nil
}

// UnmarshalU32 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalU32(x *uint32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU32 || rem < SizeHintU32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = binary.BigEndian.Uint32(buf[:])
	return buf[SizeHintU32:], rem - SizeHintU32, nil
}

// UnmarshalU64 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalU64(x *uint64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintU64 || rem < SizeHintU64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = binary.BigEndian.Uint64(buf[:])
	return buf[SizeHintU64:], rem - SizeHintU64, nil
}

// UnmarshalI8 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalI8(x *int8, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI8 || rem < SizeHintI8 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = int8(buf[0])
	return buf[SizeHintI8:], rem - SizeHintI8, nil
}

// UnmarshalI16 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalI16(x *int16, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI16 || rem < SizeHintI16 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = int16(binary.BigEndian.Uint16(buf[:]))
	return buf[SizeHintI16:], rem - SizeHintI16, nil
}

// UnmarshalI32 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalI32(x *int32, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI32 || rem < SizeHintI32 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = int32(binary.BigEndian.Uint32(buf))
	return buf[SizeHintI32:], rem - SizeHintI32, nil
}

// UnmarshalI64 from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalI64(x *int64, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintI64 || rem < SizeHintI64 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = int64(binary.BigEndian.Uint64(buf))
	return buf[SizeHintI64:], rem - SizeHintI64, nil
}
