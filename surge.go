package surge

import (
	"reflect"
	"unsafe"
)

// MaxBytes is set to 64 MB by default.
const MaxBytes = int(64 * 1024 * 1024)

// A SizeHinter can hint at the number of bytes required to represented it in
// binary.
type SizeHinter interface {
	// SizeHint returns the upper bound for the number of bytes required to
	// represent this value in binary.
	SizeHint() int
}

// A Marshaler can marshal itself into bytes.
type Marshaler interface {
	SizeHinter

	// Marshal this value into bytes.
	Marshal(buf []byte, rem int) ([]byte, int, error)
}

// An Unmarshaler can unmarshal itself from bytes.
type Unmarshaler interface {
	// Unmarshal this value from bytes.
	Unmarshal(buf []byte, rem int) ([]byte, int, error)
}

// A MarshalUnmarshaler is a marshaler and an unmarshaler.
type MarshalUnmarshaler interface {
	Marshaler
	Unmarshaler
}

// SizeHint returns the number of bytes required to store a value in its binary
// representation. This is the number of bytes "on the wire", not the number of
// bytes that need to be allocated during marshaling/unmarshaling (which can be
// different, depending on the representation of the value). SizeHint supports
// all scalars, strings, arrays, slices, maps, structs, and custom
// implementations (for types that implement the SizeHinter interface). If the
// type is not supported, then zero is returned. If the value is a pointer, then
// the size of the underlying value being pointed to will be returned.
//
//  x := int64(0)
//  sizeHint := surge.SizeHint(x)
//  if sizeHint != 8 {
//      panic("assertion failed: size of int64 must be 8 bytes")
//  }
//
func SizeHint(v interface{}) int {
	return sizeHintReflected(reflect.ValueOf(v))
}

// Marshal a value into its binary representation, and store the value in a byte
// slice. The "remaining memory quota" defines the maximum amount of bytes that
// can be allocated on the heap when marshaling the value. In this way, the
// remaining memory quota can be used to avoid allocating too much memory during
// marshaling. Marshaling supports all scalars, strings, arrays, slices, maps,
// structs, and custom implementations (for types that implement the Marshaler
// interface). After marshaling, the unconsumed tail of the byte slice, and the
// remaining memory quota, are returned. If the byte slice is too small, then an
// error is returned. Similarly, if the remaining memory quote is too small,
// then an error is returned. If the type is not supported, then an error is
// returned. An error does not imply that nothing from the byte slice, or
// remaining memory quota, was consumed. If the value is a pointer, then the
// underlying value being pointed to will be marshaled.
//
//  x := int64(0)
//  buf := make([]byte, 8)
//  tail, rem, err := surge.Marshal(x, buf, 8)
//  if len(tail) != 0 {
//      panic("assertion failed: int64 must consume 8 bytes")
//  }
//  if rem != 0 {
//      panic("assertion failed: int64 must consume 8 bytes of the memory quota")
//  }
//  if err != nil {
//      panic(fmt.Errorf("assertion failed: %v", err))
//  }
//
func Marshal(v interface{}, buf []byte, rem int) ([]byte, int, error) {
	return marshalReflected(reflect.ValueOf(v), buf, rem)
}

// Unmarshal a value from its binary representation by reading from a byte
// slice. The "remaining memory quota" defines the maximum amount of bytes that
// can be allocated on the heap when unmarshaling the value. In this way, the
// remaining memory quota can be used to avoid allocating too much memory during
// unmarshaling (this is particularly useful when dealing with potentially
// malicious input). Unmarshaling supports pointers to all scalars, strings,
// arrays, slices, maps, structs, and custom implementations (for types that
// implement the Unmarshaler interface). After unmarshaling, the unconsumed tail
// of the byte slice, and the remaining memory quota, are returned. If the byte
// slice is too small, then an error is returned. Similarly, if the remaining
// memory quote is too small, then an error is returned. If the type is not a
// pointer to one of the supported types, then an error is returned. An error
// does not imply that nothing from the byte slice, or remaining memory quota,
// was consumed. If the value is not a pointer, then an error is returned.
//
//  x := int64(0)
//  buf := make([]byte, 8)
//  tail, rem, err := surge.Unmarshal(&x, buf, 8)
//  if len(tail) != 0 {
//      panic("assertion failed: int64 must consume 8 bytes")
//  }
//  if rem != 0 {
//      panic("assertion failed: int64 must consume 8 bytes of the memory quota")
//  }
//  if err != nil {
//      panic(fmt.Errorf("assertion failed: %v", err))
//  }
//
func Unmarshal(v interface{}, buf []byte, rem int) ([]byte, int, error) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Ptr {
		return buf, rem, NewErrUnsupportedUnmarshalType(v)
	}
	return unmarshalReflected(valueOf, buf, rem)
}

func sizeHintReflected(v reflect.Value) int {
	if v.Type().Implements(sizeHinter) {
		return v.Interface().(SizeHinter).SizeHint()
	}

	switch v.Kind() {
	case reflect.Bool:
		return SizeHintBool

	case reflect.Uint8:
		return SizeHintU8
	case reflect.Uint16:
		return SizeHintU16
	case reflect.Uint32:
		return SizeHintU32
	case reflect.Uint, reflect.Uint64:
		return SizeHintU64

	case reflect.Int8:
		return SizeHintI8
	case reflect.Int16:
		return SizeHintI16
	case reflect.Int32:
		return SizeHintI32
	case reflect.Int64:
		return SizeHintI64

	case reflect.Float32:
		return SizeHintF32
	case reflect.Float64:
		return SizeHintF64

	case reflect.Array:
		return sizeHintReflectedArray(v)
	case reflect.Slice:
		return sizeHintReflectedSlice(v)
	case reflect.String:
		return sizeHintReflectedString(v)
	case reflect.Map:
		return sizeHintReflectedMap(v)
	case reflect.Struct:
		return sizeHintReflectedStruct(v)
	case reflect.Ptr:
		v = reflect.Indirect(v)
		if v.IsValid() {
			return sizeHintReflected(v)
		}
		return 0
	}

	return 0
}

func marshalReflected(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	if v.Type().Implements(marshaler) {
		return v.Interface().(Marshaler).Marshal(buf, rem)
	}

	switch v.Kind() {
	case reflect.Bool:
		return MarshalBool(v.Bool(), buf, len(buf))

	case reflect.Uint8:
		return MarshalU8(uint8(v.Uint()), buf, rem)
	case reflect.Uint16:
		return MarshalU16(uint16(v.Uint()), buf, rem)
	case reflect.Uint32:
		return MarshalU32(uint32(v.Uint()), buf, rem)
	case reflect.Uint, reflect.Uint64:
		return MarshalU64(uint64(v.Uint()), buf, rem)

	case reflect.Int8:
		return MarshalI8(int8(v.Int()), buf, rem)
	case reflect.Int16:
		return MarshalI16(int16(v.Int()), buf, rem)
	case reflect.Int32:
		return MarshalI32(int32(v.Int()), buf, rem)
	case reflect.Int64:
		return MarshalI64(int64(v.Int()), buf, rem)

	case reflect.Float32:
		return MarshalF32(float32(v.Float()), buf, rem)
	case reflect.Float64:
		return MarshalF64(float64(v.Float()), buf, rem)

	case reflect.Array:
		return marshalReflectedArray(v, buf, rem)
	case reflect.Slice:
		return marshalReflectedSlice(v, buf, rem)
	case reflect.String:
		return marshalReflectedString(v, buf, rem)
	case reflect.Map:
		return marshalReflectedMap(v, buf, rem)
	case reflect.Struct:
		return marshalReflectedStruct(v, buf, rem)
	case reflect.Ptr:
		v = reflect.Indirect(v)
		if v.IsValid() {
			return marshalReflected(v, buf, rem)
		}
		return buf, rem, nil
	}

	return buf, rem, NewErrUnsupportedMarshalType(v.Interface())
}

func unmarshalReflected(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	if v.Type().Implements(unmarshaler) {
		return v.Interface().(Unmarshaler).Unmarshal(buf, rem)
	}

	switch v.Type().Elem().Kind() {
	case reflect.Uint8:
		return UnmarshalU8((*uint8)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Uint16:
		return UnmarshalU16((*uint16)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Uint32:
		return UnmarshalU32((*uint32)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Uint64:
		return UnmarshalU64((*uint64)(unsafe.Pointer(v.Pointer())), buf, rem)

	case reflect.Int8:
		return UnmarshalI8((*int8)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Int16:
		return UnmarshalI16((*int16)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Int32:
		return UnmarshalI32((*int32)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Int64:
		return UnmarshalI64((*int64)(unsafe.Pointer(v.Pointer())), buf, rem)

	case reflect.Float32:
		return UnmarshalF32((*float32)(unsafe.Pointer(v.Pointer())), buf, rem)
	case reflect.Float64:
		return UnmarshalF64((*float64)(unsafe.Pointer(v.Pointer())), buf, rem)

	case reflect.Array:
		return unmarshalReflectedArray(v, buf, rem)
	case reflect.Slice:
		return unmarshalReflectedSlice(v, buf, rem)
	case reflect.String:
		return unmarshalReflectedString(v, buf, rem)
	case reflect.Map:
		return unmarshalReflectedMap(v, buf, rem)
	case reflect.Struct:
		return unmarshalReflectedStruct(v, buf, rem)
	}

	return buf, rem, NewErrUnsupportedUnmarshalType(v.Interface())
}

var (
	sizeHinter  = reflect.ValueOf((*SizeHinter)(nil)).Type().Elem()
	marshaler   = reflect.ValueOf((*Marshaler)(nil)).Type().Elem()
	unmarshaler = reflect.ValueOf((*Unmarshaler)(nil)).Type().Elem()
)
