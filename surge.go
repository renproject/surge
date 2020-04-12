package surge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"sort"
)

// MaxBytes is set to 8 MB by default.
const MaxBytes = int(8 * 1024 * 1024)

// Surger defines a union of the marshaler, unmarshaler, and size hinter
// interfaces.
type Surger interface {
	SizeHint() int
	Marshal(io.Writer, int) (int, error)
	Unmarshal(io.Reader, int) (int, error)
}

// SizeHinter defines an interface for types that can hint at the number of
// bytes required to represent the type in its binary form. This is useful for
// grouping memory allocations during marshaling/unmarshaling.
type SizeHinter interface {
	// SizeHint returns the recommended number of bytes that should be allocated
	// when marshaling this type. It should return an upper bound of the
	// estimated number of bytes, if the exact number is unknown.
	SizeHint() int
}

// Marshaler defines an interface for types that can be marshaled to binary.
// Marshaler types must hint at their size before marshaling.
type Marshaler interface {
	SizeHinter

	// Marshal into an I/O writer. It accepts a maximum capacity of bytes that
	// can be allocated, and returns the remaining capacity. It should not
	// allocate more bytes than the maximum capacity.
	Marshal(w io.Writer, m int) (int, error)
}

// Unmarshaler defines an interface for types that can be unmarshaled from
// binary. Unmarshaler types must hint at their size after marshaling.
type Unmarshaler interface {
	SizeHinter

	// Unmarshal from an I/O reader. It accepts a maximum capacity of bytes that
	// can be allocated, and returns the remaining capacity. It must not
	// allocate more bytes than the maximum capacity.
	Unmarshal(r io.Reader, m int) (int, error)
}

// ToBinary marshals a value into a byte slice. It allocates an in-memory
// buffer, using SizeHint to estimate the initial size of the buffer (preventing
// the need to grow the buffer while marshaling). This function is implemented
// for all scalar types, arrays, slices, and maps. If the type implements the
// Marshaler interface, then this function will use that implementation.
//
// An example of marshaling a scalar value:
//
//  x := 42
//  data, err := surge.ToBinary(x)
//  if err != nil {
//      panic(err)
//  }
//  fmt.Printf("%x", data)
//
func ToBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(SizeHint(v))
	_, err := Marshal(buf, v, MaxBytes)
	return buf.Bytes(), err
}

// FromBinary unmarshals a byte slice into a pointer. This function is
// implemented for all scalar types, arrays, slices, and maps. If the type
// implements the Unmarshaler interface, then this function will use that
// implementation. This function will never consume more memory than the default
// MaxBytes.
//
// An example of marshaling/unmarshaling a map:
//
//  xs := map[string]string{}
//  xs["foo1"] = "bar"
//  xs["foo2"] = "baz"
//
//  data, err := surge.ToBinary(xs)
//  if err != nil {
//      panic(err)
//  }
//
//  ys := map[string]string{}
//  if err := surge.FromBinary(&ys, data); err != nil {
//      panic(err)
//  }
//  fmt.Printf("foo1: %s\n", ys["foo1"])
//  fmt.Printf("foo2: %s\n", ys["foo2"])
//
func FromBinary(data []byte, v interface{}) (err error) {
	// Benchmarks indicate no immediate regression, but benchmarking
	// error cases would need to be done for deeper insights
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("recovered: %v", r)
		}
	}()

	buf := bytes.NewBuffer(data)
	_, err = Unmarshal(buf, v, MaxBytes)
	return
}

// SizeHint returns an estimate of the number of bytes that will be produced
// when marshaling a value. This is useful when pre-allocating memory to store
// marshaled values.
func SizeHint(v interface{}) int {
	switch v := v.(type) {
	case bool:
		return 1
	case *bool:
		return 1
	case []bool:
		return 4 + len(v)

	case int8:
		return 1
	case *int8:
		return 1
	case []int8:
		return 4 + len(v)

	case int16:
		return 2
	case *int16:
		return 2
	case []int16:
		return 4 + (len(v) << 1)

	case int32:
		return 4
	case *int32:
		return 4
	case []int32:
		return 4 + (len(v) << 2)

	case int64:
		return 8
	case *int64:
		return 8
	case []int64:
		return 4 + (len(v) << 3)

	case uint8:
		return 1
	case *uint8:
		return 1
	case []uint8:
		return 4 + len(v)

	case uint16:
		return 2
	case *uint16:
		return 2
	case []*uint16:
		return 4 + (len(v) << 1)

	case uint32:
		return 4
	case *uint32:
		return 4
	case []*uint32:
		return 4 + (len(v) << 2)

	case uint64:
		return 8
	case *uint64:
		return 8
	case []uint64:
		return 4 + (len(v) << 3)
	}

	if interf, ok := v.(SizeHinter); ok {
		return interf.SizeHint()
	}

	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		return SizeHint(reflect.Indirect(valOf).Interface())
	}

	switch valOf.Type().Kind() {
	case reflect.Array, reflect.Slice:
		sizeHint := 4 // Size of length prefix
		for i := 0; i < valOf.Len(); i++ {
			sizeHint += SizeHint(valOf.Index(i).Interface())
		}
		return sizeHint

	case reflect.Map:
		sizeHint := 4 // Size of length prefix
		for _, key := range valOf.MapKeys() {
			sizeHint += SizeHint(key.Interface())
			sizeHint += SizeHint(valOf.MapIndex(key).Interface())
		}
		return sizeHint
	}

	return 0
}

// Marshal a value into an I/O writer. This is more efficient than marshaling a
// value into a byte slice and returning the byte slice, because the I/O writer
// can be pre-allocated with enough memory to avoid internal allocations and
// buffer copies while marshaling.
//
// Marshaling attempts to restrict how many bytes can be allocation/written. It
// accepts the maximum number of bytes that should be allocated/written, and
// returns the number of bytes left (this can be negative).
//
// When marshaling scalars, all values are marshaled into bytes using big endian
// encoding.
//
// When marshaling arrays/slices/maps, an uint32 length prefix is marshaled and
// prefixed.
//
// When marshaling maps, key/value pairs are marshaled in order of the keys
// (sorted after the key has been marshaled). This guarantees consistency; the
// marshaled bytes are always the same if the key/values in the map are the
// same. This is particularly useful when hashing.
//
// When marshaling a value that implements the Marshaler interface, it is up to
// the user to guarantee that the implementation is sane.
func Marshal(w io.Writer, v interface{}, m int) (int, error) {
	if m <= 0 {
		return m, ErrMaxBytesExceeded
	}

	// Marshal scalar types.
	switch v := v.(type) {
	case []byte:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], uint32(len(v)))
		n, err := w.Write(bs[:])
		if err != nil {
			return m - n, err
		}
		_, err = w.Write(v)
		return m - n - len(v), err

	case string:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], uint32(len(v)))
		n, err := w.Write(bs[:])
		if err != nil {
			return m - n, err
		}
		_, err = w.Write([]byte(v))
		return m - n - len(v), err

	case bool:
		bs := [1]byte{0}
		if v {
			bs[0] = 1
		}
		n, err := w.Write(bs[:])
		return m - n, err

	case int8:
		bs := [1]byte{byte(v)}
		n, err := w.Write(bs[:])
		return m - n, err

	case int16:
		bs := [2]byte{}
		binary.BigEndian.PutUint16(bs[:], uint16(v))
		n, err := w.Write(bs[:])
		return m - n, err

	case int32:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], uint32(v))
		n, err := w.Write(bs[:])
		return m - n, err

	case int64:
		bs := [8]byte{}
		binary.BigEndian.PutUint64(bs[:], uint64(v))
		n, err := w.Write(bs[:])
		return m - n, err

	case uint8:
		bs := [1]byte{byte(v)}
		n, err := w.Write(bs[:])
		return m - n, err

	case uint16:
		bs := [2]byte{}
		binary.BigEndian.PutUint16(bs[:], v)
		n, err := w.Write(bs[:])
		return m - n, err

	case uint32:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], v)
		n, err := w.Write(bs[:])
		return m - n, err

	case uint64:
		bs := [8]byte{}
		binary.BigEndian.PutUint64(bs[:], v)
		n, err := w.Write(bs[:])
		return m - n, err
	}

	// Marshal types that implement the `Marshaler` interface.
	if interf, ok := v.(Marshaler); ok {
		return interf.Marshal(w, m)
	}

	// Marshal pointers by flattening them
	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		return Marshal(w, reflect.Indirect(valOf).Interface(), m)
	}

	// Marshal abstract data types.
	switch valOf.Type().Kind() {
	case reflect.Array, reflect.Slice:
		len := valOf.Len()
		m, err := Marshal(w, uint32(len), m)
		if err != nil {
			return m, err
		}
		for i := 0; i < len; i++ {
			m, err = Marshal(w, valOf.Index(i).Interface(), m)
			if err != nil {
				return m, err
			}
		}
		return m, nil

	case reflect.Map:
		err := error(nil)
		len := valOf.Len()

		// Sort the keys in the map, using the marshaled byte representations as
		// strings for comparison.
		marshaledKeys := make([]string, len)
		marshaledKeysToValue := make(map[string]reflect.Value, len)
		buf := new(bytes.Buffer)
		for i, key := range valOf.MapKeys() {
			// We consider the degradation of m here, so we should not consider
			// it in the next step when we write the keys to the I/O writer
			// proper.
			m, err = Marshal(buf, key.Interface(), m)
			if err != nil {
				return m, err
			}
			marshaledKey := string(buf.Bytes())
			marshaledKeys[i] = marshaledKey
			marshaledKeysToValue[marshaledKey] = valOf.MapIndex(key)
			buf.Reset()
		}
		sort.Strings(marshaledKeys)

		// Marshal the map into the writer, iterating through the sorted keys in
		// order.
		m, err := Marshal(w, uint32(len), m)
		if err != nil {
			return m, err
		}
		for _, marshaledKey := range marshaledKeys {
			// Write the key, but do not subtract the bytes written from m. We
			// have already done this in the previous step.
			_, err = w.Write([]byte(marshaledKey))
			if err != nil {
				return m, err
			}
			// Write value
			m, err = Marshal(w, marshaledKeysToValue[marshaledKey].Interface(), m)
			if err != nil {
				return m, err
			}
		}
		return m, nil
	}

	return m, newErrUnsupportedMarshalType(v)
}

// Unmarshal from an I/O reader into a pointer. This function is a complement to
// marshaling. Unmarshaling will never allocate more bytes than the given
// maximum, preventing malicious input from causing OOM errors.
func Unmarshal(r io.Reader, v interface{}, m int) (int, error) {
	if m <= 0 {
		return m, ErrMaxBytesExceeded
	}

	// Unmarshal scalar types
	switch v := v.(type) {
	case *[]byte:
		// Read length of bytes
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		len := binary.BigEndian.Uint32(bs[:])
		// Check length
		if int(len) < 0 {
			return m, newErrNegativeLength(int(len))
		}
		m -= int(len)
		if m <= 0 {
			return m, ErrMaxBytesExceeded
		}
		// Read bytes
		*v = make([]byte, len)
		_, err = io.ReadFull(r, *v)
		return m, err

	case *string:
		// Read length of string
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		len := binary.BigEndian.Uint32(bs[:])
		// Check length
		if int(len) < 0 {
			return m, newErrNegativeLength(int(len))
		}
		m -= int(len)
		if m <= 0 {
			return m, ErrMaxBytesExceeded
		}
		// Read bytes
		data := make([]byte, len)
		_, err = io.ReadFull(r, data)
		if err != nil {
			return m, err
		}
		// Read string
		*v = string(data)
		return m, nil

	case *bool:
		bs := [1]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = bs[0] != 0
		return m, nil

	case *int8:
		bs := [1]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = int8(bs[0])
		return m, nil

	case *int16:
		bs := [2]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = int16(binary.BigEndian.Uint16(bs[:]))
		return m, nil

	case *int32:
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = int32(binary.BigEndian.Uint32(bs[:]))
		return m, nil

	case *int64:
		bs := [8]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = int64(binary.BigEndian.Uint64(bs[:]))
		return m, nil

	case *uint8:
		bs := [1]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = uint8(bs[0])
		return m, nil

	case *uint16:
		bs := [2]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = binary.BigEndian.Uint16(bs[:])
		return m, nil

	case *uint32:
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = binary.BigEndian.Uint32(bs[:])
		return m, nil

	case *uint64:
		bs := [8]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		*v = binary.BigEndian.Uint64(bs[:])
		return m, nil
	}

	// Unmarshal into an interface.
	if interf, ok := v.(Unmarshaler); ok {
		return interf.Unmarshal(r, m)
	}

	// Check that we are unmarshaling into a pointer.
	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() != reflect.Ptr {
		return m, newErrUnsupportedUnmarshalType(v)
	}

	switch valOf := reflect.Indirect(valOf); valOf.Type().Kind() {
	case reflect.Array:
		// Read length of array
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		len := binary.BigEndian.Uint32(bs[:])
		// Check length
		if int(len) < 0 {
			return m, newErrNegativeLength(int(len))
		}
		if valOf.Len() != int(len) {
			return m, newErrBadLength(uint32(valOf.Len()), len)
		}
		// Read array elements
		for i := 0; i < int(len); i++ {
			m, err = Unmarshal(r, valOf.Index(i).Addr().Interface(), m)
			if err != nil {
				return m, err
			}
		}
		return m, nil

	case reflect.Slice:
		// Read length of slice
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		len := binary.BigEndian.Uint32(bs[:])
		// Check length
		if int(len) < 0 {
			return m, newErrNegativeLength(int(len))
		}

		// Scale length by the SizeHint of the element. This is because
		// each element in a list is not going to be a single byte. A similar
		// this needs to be done for maps.

		size := SizeHint(reflect.New(valOf.Type()))
		m, err = reduceBySize(m, len, size)
		if err != nil {
			return m, err
		}

		// Read slice
		valOf.Set(reflect.MakeSlice(valOf.Type(), int(len), int(len)))
		for i := 0; i < int(len); i++ {
			m, err = Unmarshal(r, valOf.Index(i).Addr().Interface(), m)
			if err != nil {
				return m, err
			}
		}
		return m, nil

	case reflect.Map:
		// Read length of map
		bs := [4]byte{}
		_, err := io.ReadFull(r, bs[:])
		if err != nil {
			return m, err
		}
		len := binary.BigEndian.Uint32(bs[:])
		// Check length
		if int(len) < 0 {
			return m, newErrNegativeLength(int(len))
		}
		// Read map
		valOf.Set(reflect.MakeMapWithSize(valOf.Type(), int(len)))
		key := reflect.New(valOf.Type().Key())
		keySize := SizeHint(key)

		m, err = reduceBySize(m, len, keySize)
		if err != nil {
			return m, err
		}

		elem := reflect.New(valOf.Type().Elem())
		elemSize := SizeHint(key)

		m, err = reduceBySize(m, len, elemSize)
		if err != nil {
			return m, err
		}

		for i := 0; i < int(len); i++ {
			// Read key
			m, err = Unmarshal(r, key.Interface(), m)
			if err != nil {
				return m, err
			}
			// Read elem
			m, err = Unmarshal(r, elem.Interface(), m)
			if err != nil {
				return m, err
			}
			// Insert into map
			valOf.SetMapIndex(reflect.Indirect(key), reflect.Indirect(elem))
		}
		return m, nil
	}

	return m, newErrUnsupportedUnmarshalType(v)
}

func reduceBySize(m int, len uint32, size int) (int, error) {
	if size <= 16 {
		size = 16 // All unknown types will take up approximately 16 bytes, assuming a 64-bit system.
	}
	if int(len)*size < int(len) {
		return m, newErrLengthOverflow()
	}
	if int(len)*size < 0 {
		return m, newErrLengthOverflow()
	}
	if int(len)*size > m {
		return m, newErrLengthOverflow()
	}

	m -= int(len) * size

	if m <= 0 {
		return m, ErrMaxBytesExceeded
	}
	return m, nil
}
