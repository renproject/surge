package surge

// MarshalLen marshalls the given slice length.
func MarshalLen(l uint32, buf []byte, rem int) ([]byte, int, error) {
	return MarshalU32(l, buf, rem)
}

// UnmarshalLen unmarshalls a slice length, checking that the total space
// required for the slice will not exceed rem.
//
// Panics: This function will panic if elemSize is 0.
func UnmarshalLen(dst *uint32, elemSize int, buf []byte, rem int) ([]byte, int, error) {
	var l uint32
	buf, rem, err := UnmarshalU32(&l, buf, rem)
	if err != nil {
		return buf, rem, err
	}

	var c uint64 = uint64(l) * uint64(elemSize)

	// Check if there was overflow in the multiplication.
	// NOTE: This can only occur on systems for which the int type is 64 bits,
	// and in addition when elemSize >= 2^32 = 4 Gb. Elements with this size
	// are unlikely to be used in practice.
	if c/uint64(elemSize) != uint64(l) {
		return buf, rem, ErrLengthOverflow
	}

	if uint64(rem) < c {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}

	*dst = l
	return buf, rem, nil
}
