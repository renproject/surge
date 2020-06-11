package surge

import (
	"reflect"
	"sort"
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
	n := v.Len() * int(v.Type().Elem().Size())
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
		if _, rem, err = marshalReflected(key, keyValue.keyData, rem); err != nil {
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
			return buf, rem, ErrMaxBytesExceeded
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
	mapLen := uint16(0)
	buf, rem, err := UnmarshalU16(&mapLen, buf, rem)
	if err != nil {
		return buf, rem, err
	}

	n := int(mapLen)
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	n *= int(v.Type().Elem().Size())
	if n < 0 {
		return buf, rem, ErrLengthOverflow
	}
	if rem < n {
		return buf, rem, ErrUnexpectedEndOfBuffer
	}
	rem -= n
	v.Set(reflect.MakeMapWithSize(v.Type(), int(mapLen)))

	for i := uint16(0); i < mapLen; i++ {
		key := reflect.New(v.Type().Key())
		elem := reflect.New(v.Type().Elem())
		if buf, rem, err = unmarshalReflected(key, buf, rem); err != nil {
			return buf, rem, err
		}
		if buf, rem, err = unmarshalReflected(elem, buf, rem); err != nil {
			return buf, rem, err
		}
		v.SetMapIndex(reflect.Indirect(key), reflect.Indirect(elem))
	}
	return buf, rem, nil
}
