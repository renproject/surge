package surge

import (
	"io"
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

// A Sinker can stream itself into an IO writer. It will have the same format as
// when it is marshaling, but it will write directly into the IO writer. This is
// particularly useful for writing into size-limited IO buffers.
type Sinker interface {
	SizeHinter

	// Sink the bytes of this value into an IO writer.
	Sink(io.Writer) error
}

// A Streamer can stream from an IO reader into itself. It will expect the same
// format as when it is unmarshaling, but it will read directly from the IO
// reader. This is particularly useful for reading from size-limited IO buffers.
type Streamer interface {
	// Stream bytes from an IO reader into this value.
	Stream(io.Reader) error
}

// A SinkStreamer is a sinker and a streamer.
type SinkStreamer interface {
	Sinker
	Streamer
}

// A MarshalUnmarshaler is a marshaler and an unmarshaler.
type MarshalUnmarshaler interface {
	Marshaler
	Unmarshaler
}

func SizeHint(v interface{}) int {
	return sizeHintReflected(reflect.ValueOf(v))
}

func Marshal(v interface{}, buf []byte, rem int) ([]byte, int, error) {
	return marshalReflected(reflect.ValueOf(v), buf, rem)
}

func Unmarshal(v interface{}, buf []byte, rem int) ([]byte, int, error) {
	valueOf := reflect.ValueOf(v)
	if valueOf.Kind() != reflect.Ptr {
		return buf, rem, newErrUnsupportedUnmarshalType(v)
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
	case reflect.Int, reflect.Int64:
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
	case reflect.Int, reflect.Int64:
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

	return buf, rem, newErrUnsupportedMarshalType(v.Interface())
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

	return buf, rem, newErrUnsupportedUnmarshalType(v.Interface())
}

var (
	sizeHinter  = reflect.ValueOf((*SizeHinter)(nil)).Type().Elem()
	marshaler   = reflect.ValueOf((*Marshaler)(nil)).Type().Elem()
	unmarshaler = reflect.ValueOf((*Unmarshaler)(nil)).Type().Elem()
)
