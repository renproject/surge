package surgeutil

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing/quick"
	"time"

	"github.com/renproject/surge"
)

// MarshalUnmarshalCheck generates a random instance of a type, marshals it into
// binary, unmarshals the result into a new instance of the type, and then does
// a deep equality check on the two instances. An error is returned when
// generation, marshaling, or unmarshaling return an error, or if the two
// instances are unequal. Otherwise, it returns nil.
func MarshalUnmarshalCheck(t reflect.Type) error {
	// Generate
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal
	data, err := surge.ToBinary(x.Interface())
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}
	// Unmarshal
	y := reflect.New(t)
	if err := surge.FromBinary(y.Interface(), data); err != nil {
		return fmt.Errorf("cannot unmarshal: %v", err)
	}
	// Equality
	if !reflect.DeepEqual(x.Interface(), y.Elem().Interface()) {
		return fmt.Errorf("unequal")
	}
	return nil
}

// Fuzz generates a random instance of a type, and then attempts to unmarshal
// random bytes into that instance. It returns nothing, but is expected to not
// panic.
func Fuzz(t reflect.Type) {
	// Fuzz data
	data, ok := quick.Value(reflect.TypeOf([]byte{}), rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		panic(fmt.Errorf("cannot generate value of type %v", t))
	}
	// Unmarshal
	x := reflect.New(t)
	if err := surge.FromBinary(x.Interface(), data.Bytes()); err != nil {
		// Ignore the error, because we are only interested in whether or not
		// the unmarshaling causes a panic.
	}
}

// MarshalBufTooSmallSparse generates a random intance of a type, and then
// attempts to marshal it into a buffer that is too small. It returns an error
// when marshaling succeeds. Otherwise, it returns nil.
//
// The number of buffer sizes tested is determinted by the given integer. A
// value of 0 means that all buffer sizes will be tested.
func MarshalBufTooSmallSparse(t reflect.Type, steps int) error {
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}

	size := surge.SizeHint(x.Interface())
	step := stepSize(size, steps)
	fullBuf := make([]byte, size)
	for bufLen := 0; bufLen < size; bufLen += step {
		buf := fullBuf[:bufLen]
		rem := size
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected success: %v < %v", bufLen, size)
		}
	}
	return nil
}

// MarshalBufTooSmall generates a random intance of a type, and then attempts to
// marshal it into a buffer that is too small. It returns an error when
// marshaling succeeds. Otherwise, it returns nil.
//
// Equivalent to MarshalBufTooSmallSparse(t, 0).
func MarshalBufTooSmall(t reflect.Type) error {
	return MarshalBufTooSmallSparse(t, 0)
}

// MarshalRemTooSmallSparse generates a random intance of a type, and then
// attempts to marshal when a remaining memory quota that is too small. It
// returns an error when marshaling succeeds. Otherwise, it returns nil.
//
// The number of buffer sizes tested is determinted by the given integer. A
// value of 0 means that all buffer sizes will be tested.
func MarshalRemTooSmallSparse(t reflect.Type, steps int) error {
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}

	size := surge.SizeHint(x.Interface())
	step := stepSize(size, steps)
	buf := make([]byte, size)
	for rem := 0; rem < size; rem += step {
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

// MarshalRemTooSmall generates a random intance of a type, and then attempts to
// marshal when a remaining memory quota that is too small. It returns an error
// when marshaling succeeds. Otherwise, it returns nil.
//
// Equivalent to MarshalRemTooSmallSparse(t, 0).
func MarshalRemTooSmall(t reflect.Type) error {
	return MarshalRemTooSmallSparse(t, 0)
}

// UnmarshalBufTooSmallSparse generates a random intance of a type, marshals it
// into binary, and then attempts to unmarshal the result with a buffer that is
// too small. It returns an error when marshaling succeeds. Otherwise, it
// returns nil.
//
// The number of buffer sizes tested is determinted by the given integer. A
// value of 0 means that all buffer sizes will be tested.
func UnmarshalBufTooSmallSparse(t reflect.Type, steps int) error {
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}

	buf, err := surge.ToBinary(x.Interface())
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}

	step := stepSize(len(buf), steps)
	y := reflect.New(t)
	for bufLen := 0; bufLen < len(buf); bufLen += step {
		if _, _, err := surge.Unmarshal(y.Interface(), buf[:bufLen], surge.MaxBytes); err == nil {
			return fmt.Errorf("unexpected success: %v < %v", bufLen, len(buf))
		}
	}
	return nil
}

// UnmarshalBufTooSmall generates a random intance of a type, marshals it into
// binary, and then attempts to unmarshal the result with a buffer that is too
// small. It returns an error when marshaling succeeds. Otherwise, it returns
// nil.
//
// Equivalent to UnmarshalBufTooSmallSparse(t, 0).
func UnmarshalBufTooSmall(t reflect.Type) error {
	return UnmarshalBufTooSmallSparse(t, 0)
}

// UnmarshalRemTooSmallSparse generates a random intance of a type, marshals it
// into binary, and then attempts to unmarshal the result with a remaining
// memory quota that is too small. It returns an error when marshaling
// succeeds.  Otherwise, it returns nil.
//
// The number of buffer sizes tested is determinted by the given integer. A
// value of 0 means that all buffer sizes will be tested.
func UnmarshalRemTooSmallSparse(t reflect.Type, steps int) error {
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}

	size := surge.SizeHint(x.Interface())
	buf := make([]byte, size)
	if _, _, err := surge.Marshal(x.Interface(), buf, surge.MaxBytes); err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}

	rem := size
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when unmarshaling
		rem += x.Len() * int(t.Key().Size()+t.Elem().Size())
	}

	step := stepSize(rem, steps)
	y := reflect.New(t)
	for rem2 := 0; rem2 < rem; rem2 += step {
		if _, _, err := surge.Unmarshal(y.Interface(), buf, rem2); err == nil {
			return fmt.Errorf("unexpected success: %v < %v", rem2, rem)
		}
	}
	return nil
}

// UnmarshalRemTooSmall generates a random intance of a type, marshals it into
// binary, and then attempts to unmarshal the result with a remaining memory
// quota that is too small. It returns an error when marshaling succeeds.
// Otherwise, it returns nil.
//
// Equivalent to UnmarshalRemTooSmallSparse(t, 0).
func UnmarshalRemTooSmall(t reflect.Type) error {
	return UnmarshalRemTooSmallSparse(t, 0)
}

func stepSize(max, steps int) int {
	var step int
	if steps == 0 {
		step = 1
	} else {
		step = max / steps
		if step < 1 {
			step = 1
		}
	}
	return step
}
