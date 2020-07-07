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

// MarshalBufTooSmall generates a random intance of a type, and then attempts to
// marshal it into a buffer that is too small. It returns an error when
// marshaling succeeds. Otherwise, it returns nil.
func MarshalBufTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Generate buffer that is too small
	size := surge.SizeHint(x.Interface())
	// Marshal
	for bufLen := 0; bufLen < size; bufLen++ {
		buf := make([]byte, bufLen)
		rem := size
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

// MarshalRemTooSmall generates a random intance of a type, and then attempts to
// marshal when a remaining memory quota that is too small. It returns an error
// when marshaling succeeds. Otherwise, it returns nil.
func MarshalRemTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Generate buffer that is too small
	size := surge.SizeHint(x.Interface())
	// Marshal
	for rem := 0; rem < size; rem++ {
		buf := make([]byte, size)
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

// UnmarshalBufTooSmall generates a random intance of a type, marshals it into
// binary, and then attempts to unmarshal the result with a buffer that is too
// small. It returns an error when marshaling succeeds. Otherwise, it returns
// nil.
func UnmarshalBufTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal the value so that we can attempt to unmarshal the resulting data
	buf, err := surge.ToBinary(x.Interface())
	if err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	// Unmarshal with buffers that are too small
	for bufLen := 0; bufLen < len(buf); bufLen++ {
		y := reflect.New(t)
		if _, _, err := surge.Unmarshal(y.Interface(), buf[:bufLen], surge.MaxBytes); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

// UnmarshalRemTooSmall generates a random intance of a type, marshals it into
// binary, and then attempts to unmarshal the result with a remaining memory
// quota that is too small. It returns an error when marshaling succeeds.
// Otherwise, it returns nil.
func UnmarshalRemTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal the value so that we can attempt to unmarshal the resulting data
	size := surge.SizeHint(x.Interface())
	buf := make([]byte, size)
	rem := size
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when marshaling
		rem = size + 48*x.Len()
	}
	if _, _, err := surge.Marshal(x.Interface(), buf, rem); err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	// Unmarshal with memory quotas that are too small
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when unmarshaling
		rem = size + x.Len()*int(t.Key().Size()+t.Elem().Size())
	}
	for rem2 := 0; rem2 < rem; rem2++ {
		y := reflect.New(t)
		if _, _, err := surge.Unmarshal(y.Interface(), buf, rem2); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}
