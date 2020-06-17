package surge

import (
	"reflect"
	"sort"
	"unsafe"
)

func sizeHintReflectedMap(v reflect.Value) int {
	sizeHint := 2
	iter := v.MapRange()
	for iter.Next() {
		sizeHint += sizeHintReflected(iter.Key())
		sizeHint += sizeHintReflected(iter.Value())
	}
	return sizeHint
}

func marshalReflectedMap(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	buf, rem, err := MarshalU16(uint16(v.Len()), buf, rem)
	if err != nil {
		return buf, rem, err
	}

	// Define a local key/value type that can be used to sort the key/value
	// pairs. It is worth noting, the key is explicitly defined to be bytes, so
	// we will need to marshal all keys and then convert those bytes to strings.
	// This guarantees that, regardless of the key type, we can sort the keys
	// deterministically.
	type KeyValue struct {
		keyData []byte
		key     reflect.Value
	}

	// Allocate a slice for storing the key/value pairs. This is needed so that
	// we can then sort the slice, ensuring that the key/value ordering is
	// deterministic.
	n := v.Len() * int(unsafe.Sizeof(KeyValue{}))
	if n < 0 || n < v.Len() {
		return buf, rem, ErrLengthOverflow
	}
	if rem < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	rem -= n
	keyValues := make([]KeyValue, 0, v.Len())

	for _, key := range v.MapKeys() {
		// Marshal the key into bytes, so that we can guarantee that the keys
		// can be compared.
		sizeHint := sizeHintReflected(key)
		if rem < sizeHint {
			return buf, rem, ErrUnexpectedEndOfBuffer
		}
		rem -= sizeHint
		keyValue := KeyValue{
			keyData: make([]byte, sizeHint),
			key:     key,
		}
		if _, rem, err = marshalReflected(key, keyValue.keyData, rem+sizeHint); err != nil {
			return buf, rem, err
		}

		// Search and insert to ensure that the key/values are always in sorted
		// order.
		i := sort.Search(len(keyValues), func(i int) bool {
			other := keyValues[i]
			if len(keyValue.keyData) < len(other.keyData) {
				return true
			}
			if len(keyValue.keyData) > len(other.keyData) {
				return false
			}
			for j := range keyValue.keyData {
				if keyValue.keyData[j] < other.keyData[j] {
					return true
				}
				if keyValue.keyData[j] > other.keyData[j] {
					return false
				}
			}
			return true
		})
		keyValues = append(keyValues, KeyValue{})
		copy(keyValues[i+1:], keyValues[i:])
		keyValues[i] = keyValue
	}

	for _, keyValue := range keyValues {
		// Marshal the key. This is the second time we are doing it, so we do
		// not check/decrement rem again. This is the first time we are actually
		// writing out to the buffer though, so we do need to check buffer
		// lengths.
		keyDataLen := len(keyValue.keyData)
		if len(buf) < keyDataLen {
			return buf, rem, ErrUnexpectedEndOfBuffer
		}
		copy(buf, keyValue.keyData)
		buf = buf[keyDataLen:]

		// Marshal the value.
		if buf, rem, err = marshalReflected(v.MapIndex(keyValue.key), buf, rem); err != nil {
			return buf, rem, err
		}
	}
	return buf, rem, nil
}

func unmarshalReflectedMap(v reflect.Value, buf []byte, rem int) ([]byte, int, error) {
	var err error

	mapLen := uint16(0)
	if buf, rem, err = UnmarshalU16(&mapLen, buf, rem); err != nil {
		return buf, rem, err
	}

	elem := v.Elem()
	n := int(mapLen)
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	size := int(elem.Type().Size())
	if rem < size || size < 0 {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	rem -= size
	elem.Set(reflect.MakeMapWithSize(elem.Type(), int(mapLen)))

	for i := uint16(0); i < mapLen; i++ {
		k := reflect.New(elem.Type().Key())
		v := reflect.New(elem.Type().Elem())
		if buf, rem, err = unmarshalReflected(k, buf, rem); err != nil {
			return buf, rem, err
		}
		if buf, rem, err = unmarshalReflected(v, buf, rem); err != nil {
			return buf, rem, err
		}
		elem.SetMapIndex(reflect.Indirect(k), reflect.Indirect(v))
	}
	return buf, rem, nil
}
