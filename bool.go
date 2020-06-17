package surge

const (
	// SizeHintBool is the number of bytes required to represent a boolean value
	// in binary.
	SizeHintBool = 1
)

// MarshalBool into a byte slice. It will not consume more memory than the
// remaining memory quota (either through writes, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func MarshalBool(x bool, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintBool || rem < SizeHintBool {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	if x {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf[SizeHintBool:], rem - SizeHintBool, nil
}

// UnmarshalBool from a byte slice. It will not consume more memory than the
// remaining memory quota (either through reads, or in-memory allocations). It
// will return the unconsumed tail of the byte slice, and the remaining memory
// quota. An error is returned if the byte slice is too small, or if the
// remainin memory quote is insufficient.
func UnmarshalBool(x *bool, buf []byte, rem int) ([]byte, int, error) {
	if len(buf) < SizeHintBool || rem < SizeHintBool {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	*x = buf[0] != 0
	return buf[SizeHintBool:], rem - SizeHintBool, nil
}
