package surge

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"sort"
)

// MaxBytes is set to 8 MB by default.
var MaxBytes = uint32(8 * 1024 * 1024)

// Marshaler defines an interface for types that can be marshaled to binary.
// Marshaler types must hint at their size before marshaling, and must return
// the number of bytes written after marshaling.
type Marshaler interface {
	SizeHinter
	Marshal(w io.Writer) (n uint32, err error)
}

// Unmarshaler defines an interface for types that can be unmarshaled from
// binary. Unmarshaler types must hint at their size after marshaling, and must
// return the number of bytes read after unmarshaling. Unmarshaler types must
// also limit the size of memory allocations during unmarshaling.
type Unmarshaler interface {
	SizeHinter
	Unmarshal(r io.Reader, m uint32) (n uint32, err error)
}

// SizeHinter defines an interface for types that can hint at the number of
// bytes required to represent the type in its binary form. This is useful for
// grouping memory allocations during marshaling/unmarshaling.
type SizeHinter interface {
	SizeHint() (n uint32)
}

// ToBinary marshals a value into a byte slice. It allocates an in-memory
// buffer, using SizeHint to estimate the initial size of the buffer (preventing
// the need to grow the buffer during marshaling). This function should not be
// used to implement the Marshaler interface, or as part of marshaling a parent
// value.
//
// An example of marshaling a scalar value:
//
//	x := 42
//	data, err := surge.ToBinary(x)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%x", data)
//
// An example of marshaling a custom struct:
//
//	type Point struct {
//		X uint64 `surge:"0"`
//		Y uint64 `surge:"1"`
//	}
//
//	p := Point{ 13, 169 }
//	data, err := surge.ToBinary(p)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Printf("%x", data)
//
func ToBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(int(SizeHint(v)))
	_, err := Marshal(v, buf)
	return buf.Bytes(), err
}

// FromBinary unmarshals a byte slice into a destination value. The destination
// value must be a pointer.
//
// An example of marshaling/unmarshaling a map:
//
//	xs := map[string]string{}
//	xs["foo1"] = "bar"
//	xs["foo2"] = "baz"
//
//	data, err := surge.ToBinary(xs)
//	if err != nil {
//		panic(err)
//	}
//
//	ys := map[string]string{}
//	if err := surge.FromBinary(&ys, data); err != nil {
//		panic(err)
//	}
//	fmt.Printf("foo1: %s\n", ys["foo1"])
//	fmt.Printf("foo2: %s\n", ys["foo2"])
//
func FromBinary(v interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	_, err := Unmarshal(v, buf, MaxBytes)
	return err
}

// Marshal a value into an io.Writer. This is more efficient than marshaling a
// value into a byte slice and returning the byte slice, because the io.Writer
// can be pre-allocated with enough memory to avoid internal allocations and
// buffer copies. This function should only used when defining custom
// implementations of the Marshaler interface, or as part of marshaling a parent
// value. In most use cases, the ToBinary function should be used.
//
// When marshaling scalars, all values are marshaled into bytes using little
// endian encoding.
//
// When marshaling arrays/slices/maps, an uint32 length prefix is marshaled and
// prefixed.
//
// When marshaling maps, key/value pairs are marshaled in order of the keys
// (sorted after the key has been marshaled). This guarantees consistency; the
// marshaled bytes are always the same if the key/values in the map are the
// same. This is particularly useful when hashing.
//
// When marshaling custom struct, struct tags are used to convert the struct
// into a map (which is then marshaled like a normal map).
//
// When marshaling a value that implements the Marshaler interface, it is up to
// the user to guarantee that the implementation is sane.
func Marshal(v interface{}, w io.Writer) (uint32, error) {
	// Marshal scalar types.
	switch v := v.(type) {
	case []byte:
		len := len(v)
		n1, err := Marshal(uint32(len), w)
		if err != nil {
			return n1, err
		}
		n2, err := w.Write(v)
		return n1 + uint32(n2), err

	case bool:
		bs := [1]byte{0}
		if v {
			bs[0] = 1
		}
		n, err := w.Write(bs[:])
		return uint32(n), err

	case int8:
		bs := [1]byte{byte(v)}
		n, err := w.Write(bs[:])
		return uint32(n), err

	case int16:
		bs := [2]byte{}
		binary.BigEndian.PutUint16(bs[:], uint16(v))
		n, err := w.Write(bs[:])
		return uint32(n), err

	case int32:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], uint32(v))
		n, err := w.Write(bs[:])
		return uint32(n), err

	case int64:
		bs := [8]byte{}
		binary.BigEndian.PutUint64(bs[:], uint64(v))
		n, err := w.Write(bs[:])
		return uint32(n), err

	case uint8:
		bs := [1]byte{byte(v)}
		n, err := w.Write(bs[:])
		return uint32(n), err

	case uint16:
		bs := [2]byte{}
		binary.BigEndian.PutUint16(bs[:], v)
		n, err := w.Write(bs[:])
		return uint32(n), err

	case uint32:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], v)
		n, err := w.Write(bs[:])
		return uint32(n), err

	case uint64:
		bs := [8]byte{}
		binary.BigEndian.PutUint64(bs[:], v)
		n, err := w.Write(bs[:])
		return uint32(n), err

	case string:
		bs := [4]byte{}
		binary.BigEndian.PutUint32(bs[:], uint32(len(v)))
		n1, err := w.Write(bs[:])
		if err != nil {
			return uint32(n1), err
		}
		n2, err := w.Write([]byte(v))
		return uint32(n1 + n2), err
	}

	// Marshal types that implement the `Marshaler` interface.
	if i, ok := v.(Marshaler); ok {
		return i.Marshal(w)
	}

	// Marshal pointers by flattening them
	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		return Marshal(reflect.Indirect(valOf).Interface(), w)
	}

	// Marshal abstract data types.
	switch valOf.Type().Kind() {
	case reflect.Array, reflect.Slice:
		len := valOf.Len()
		n1, err := Marshal(uint32(len), w)
		if err != nil {
			return n1, err
		}
		for i := 0; i < len; i++ {
			n2, err := Marshal(valOf.Index(i).Interface(), w)
			n1 += n2
			if err != nil {
				return n1, err
			}
		}
		return n1, nil

	case reflect.Map:
		len := valOf.Len()

		// Sort the keys in the map, using the marshaled byte representations as
		// strings for comparison.
		marshaledKeys := make([]string, len)
		marshaledKeysToValue := make(map[string]reflect.Value, len)
		buf := new(bytes.Buffer)
		for i, key := range valOf.MapKeys() {
			if _, err := Marshal(key.Interface(), buf); err != nil {
				return 0, err
			}
			marshaledKey := string(buf.Bytes())
			marshaledKeys[i] = marshaledKey
			marshaledKeysToValue[marshaledKey] = valOf.MapIndex(key)
			buf.Reset()
		}
		sort.Strings(marshaledKeys)

		// Marshal the map into the writer, iterating through the sorted keys in
		// order.
		n1, err := Marshal(uint32(len), w)
		if err != nil {
			return n1, err
		}
		for _, marshaledKey := range marshaledKeys {
			n2, err := w.Write([]byte(marshaledKey))
			n1 += uint32(n2)
			if err != nil {
				return n1, err
			}
			n3, err := Marshal(marshaledKeysToValue[marshaledKey].Interface(), w)
			n1 += n3
			if err != nil {
				return n1, err
			}
		}
		return n1, nil
	}

	return 0, newErrUnsupportedMarshalType(v)
}

// Unmarshal from an io.Reader into a destination value. The destination value
// must be a pointer. This function should only used when defining custom
// implementations of the Unmarshaler interface, or as part of marshaling a
// parent value. In most use cases, the FromBinary function should be used.
func Unmarshal(v interface{}, r io.Reader, m uint32) (uint32, error) {
	switch v := v.(type) {
	case *[]byte:
		bs := [4]byte{}
		n1, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n1), err
		}
		len := binary.BigEndian.Uint32(bs[:])
		if len > m {
			return uint32(n1), ErrMaxBytesExceeded
		}
		*v = make([]byte, len)
		n2, err := io.ReadFull(r, *v)
		if err != nil {
			return uint32(n1 + n2), err
		}
		return uint32(n1 + n2), nil

	case *bool:
		bs := [1]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = bs[0] != 0
		return uint32(n), nil

	case *int8:
		bs := [1]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = int8(bs[0])
		return uint32(n), nil

	case *int16:
		bs := [2]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = int16(binary.BigEndian.Uint16(bs[:]))
		return uint32(n), nil

	case *int32:
		bs := [4]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = int32(binary.BigEndian.Uint32(bs[:]))
		return uint32(n), nil

	case *int64:
		bs := [8]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = int64(binary.BigEndian.Uint64(bs[:]))
		return uint32(n), nil

	case *uint8:
		bs := [1]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = uint8(bs[0])
		return uint32(n), nil

	case *uint16:
		bs := [2]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = binary.BigEndian.Uint16(bs[:])
		return uint32(n), nil

	case *uint32:
		bs := [4]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = binary.BigEndian.Uint32(bs[:])
		return uint32(n), nil

	case *uint64:
		bs := [8]byte{}
		n, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n), err
		}
		*v = binary.BigEndian.Uint64(bs[:])
		return uint32(n), nil

	case *string:
		bs := [4]byte{}
		n1, err := io.ReadFull(r, bs[:])
		if err != nil {
			return uint32(n1), err
		}
		len := binary.BigEndian.Uint32(bs[:])
		if len > m {
			return uint32(n1), ErrMaxBytesExceeded
		}
		data := make([]byte, len)
		n2, err := io.ReadFull(r, data)
		if err != nil {
			return uint32(n1 + n2), err
		}
		*v = string(data)
		return uint32(n1 + n2), nil
	}

	if i, ok := v.(Unmarshaler); ok {
		return i.Unmarshal(r, m)
	}

	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		switch valOf := reflect.Indirect(valOf); valOf.Type().Kind() {
		case reflect.Array:
			len := uint32(0)
			n1, err := Unmarshal(&len, r, m)
			if err != nil {
				return n1, err
			}
			m -= n1 + len
			if m < 0 {
				return n1, ErrMaxBytesExceeded
			}
			if uint32(valOf.Len()) != len {
				return n1, newErrBadLength(uint32(valOf.Len()), len)
			}
			for i := 0; i < int(len); i++ {
				n2, err := Unmarshal(valOf.Index(i).Addr().Interface(), r, m)
				n1 += n2
				if err != nil {
					return n1, err
				}
				m -= n2
				if m < 0 {
					return n1, ErrMaxBytesExceeded
				}
			}
			return n1, nil

		case reflect.Slice:
			len := uint32(0)
			n1, err := Unmarshal(&len, r, m)
			if err != nil {
				return n1, err
			}
			m -= n1 + len
			if m < 0 {
				return n1, ErrMaxBytesExceeded
			}
			valOf.Set(reflect.MakeSlice(valOf.Type(), int(len), int(len)))
			for i := 0; i < int(len); i++ {
				n2, err := Unmarshal(valOf.Index(i).Addr().Interface(), r, m)
				n1 += n2
				if err != nil {
					return n1, err
				}
				m -= n2
				if m < 0 {
					return n1, ErrMaxBytesExceeded
				}
			}
			return n1, nil

		case reflect.Map:
			len := uint32(0)
			n1, err := Unmarshal(&len, r, m)
			if err != nil {
				return n1, err
			}
			m -= n1 + len
			if m < 0 {
				return n1, ErrMaxBytesExceeded
			}
			valOf.Set(reflect.MakeMapWithSize(valOf.Type(), int(len)))
			key := reflect.New(valOf.Type().Key())
			elem := reflect.New(valOf.Type().Elem())
			for i := 0; i < int(len); i++ {
				n2, err := Unmarshal(key.Interface(), r, m)
				n1 += n2
				if err != nil {
					return n1, err
				}
				m -= n2
				if m < 0 {
					return n1, ErrMaxBytesExceeded
				}

				n3, err := Unmarshal(elem.Interface(), r, m)
				n1 += n3
				if err != nil {
					return n1, err
				}
				m -= n3
				if m < 0 {
					return n1, ErrMaxBytesExceeded
				}

				valOf.SetMapIndex(reflect.Indirect(key), reflect.Indirect(elem))
			}
			return n1, nil
		}
	}

	return 0, newErrUnsupportedUnmarshalType(v)
}

// SizeHint returns an estimate of the number of bytes that will be produced
// when marshaling a value. This is useful when pre-allocating memory to store
// marshaled values.
func SizeHint(v interface{}) uint32 {
	switch v := v.(type) {
	case bool:
		return 1
	case *bool:
		return 1
	case []bool:
		return 4 + uint32(len(v))

	case int8:
		return 1
	case *int8:
		return 1
	case []int8:
		return 4 + uint32(len(v))

	case int16:
		return 2
	case *int16:
		return 2
	case []int16:
		return 4 + uint32(len(v)<<1)

	case int32:
		return 4
	case *int32:
		return 4
	case []int32:
		return 4 + uint32(len(v)<<2)

	case int64:
		return 8
	case *int64:
		return 8
	case []int64:
		return 4 + uint32(len(v)<<3)

	case uint8:
		return 1
	case *uint8:
		return 1
	case []uint8:
		return 4 + uint32(len(v))

	case uint16:
		return 2
	case *uint16:
		return 2
	case []*uint16:
		return 4 + uint32(len(v)<<1)

	case uint32:
		return 4
	case *uint32:
		return 4
	case []*uint32:
		return 4 + uint32(len(v)<<2)

	case uint64:
		return 8
	case *uint64:
		return 8
	case []uint64:
		return 4 + uint32(len(v)<<3)
	}

	if i, ok := v.(SizeHinter); ok {
		return i.SizeHint()
	}

	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		return SizeHint(reflect.Indirect(valOf).Interface())
	}

	switch valOf.Type().Kind() {
	case reflect.Array, reflect.Slice:
		sizeHint := uint32(4) // Size of length prefix
		for i := 0; i < valOf.Len(); i++ {
			sizeHint += SizeHint(valOf.Index(i).Interface())
		}
		return sizeHint

	case reflect.Map:
		sizeHint := uint32(4) // Size of length prefix
		for _, key := range valOf.MapKeys() {
			sizeHint += SizeHint(key.Interface())
			sizeHint += SizeHint(valOf.MapIndex(key).Interface())
		}
		return sizeHint
	}

	return 0
}
