package surge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
)

type Marshaler interface {
	SizeHinter
	Marshal(io.Writer) error
}

type Unmarshaler interface {
	SizeHinter
	Unmarshal(io.Reader) error
}

type SizeHinter interface {
	SizeHint() int
}

func ToBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(SizeHint(v))
	err := Marshal(v, buf)
	return buf.Bytes(), err
}

func FromBinary(v interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	return Unmarshal(v, buf)
}

func Marshal(v interface{}, w io.Writer) error {
	// Marshal primitive types.
	switch v := v.(type) {
	case bool:
		bs := [1]byte{0}
		if v {
			bs[0] = 1
		}
		_, err := w.Write(bs[:])
		return err

	case int8:
		bs := [1]byte{byte(v)}
		_, err := w.Write(bs[:])
		return err

	case int16:
		bs := [2]byte{}
		binary.LittleEndian.PutUint16(bs[:], uint16(v))
		_, err := w.Write(bs[:])
		return err

	case int32:
		bs := [4]byte{}
		binary.LittleEndian.PutUint32(bs[:], uint32(v))
		_, err := w.Write(bs[:])
		return err

	case int64:
		bs := [8]byte{}
		binary.LittleEndian.PutUint64(bs[:], uint64(v))
		_, err := w.Write(bs[:])
		return err

	case uint8:
		bs := [1]byte{byte(v)}
		_, err := w.Write(bs[:])
		return err

	case uint16:
		bs := [2]byte{}
		binary.LittleEndian.PutUint16(bs[:], v)
		_, err := w.Write(bs[:])
		return err

	case uint32:
		bs := [4]byte{}
		binary.LittleEndian.PutUint32(bs[:], v)
		_, err := w.Write(bs[:])
		return err

	case uint64:
		bs := [8]byte{}
		binary.LittleEndian.PutUint64(bs[:], v)
		_, err := w.Write(bs[:])
		return err

	case string:
		bs := [4]byte{}
		binary.LittleEndian.PutUint32(bs[:], uint32(len(v)))
		if _, err := w.Write(bs[:]); err != nil {
			return err
		}
		_, err := w.Write([]byte(v))
		return err
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
		if err := Marshal(uint32(len), w); err != nil {
			return err
		}
		for i := 0; i < len; i++ {
			if err := Marshal(valOf.Index(i).Interface(), w); err != nil {
				return err
			}
		}
		return nil

	case reflect.Map:
		len := valOf.Len()
		keys := make([]string, len)
		keysLookup := make(map[string]reflect.Value, len)
		buf := new(bytes.Buffer)
		for i, key := range valOf.MapKeys() {
			if err := Marshal(key.Interface(), buf); err != nil {
				return err
			}
			marshaledKeyAsStr := string(buf.Bytes())
			keys[i] = marshaledKeyAsStr
			keysLookup[marshaledKeyAsStr] = key
			buf.Reset()
		}
		sort.Strings(keys)
		if err := Marshal(uint32(len), w); err != nil {
			return err
		}
		for _, key := range keys {
			if err := Marshal(keysLookup[key].Interface(), w); err != nil {
				return err
			}
			if err := Marshal(valOf.MapIndex(keysLookup[key]).Interface(), w); err != nil {
				return err
			}
		}
		return nil

	case reflect.Struct:
		adt := make(map[uint16]interface{}, valOf.Type().NumField())
		for i := 0; i < valOf.Type().NumField(); i++ {
			field := valOf.Type().Field(i)
			tagString := field.Tag.Get("surge")
			if tagString == "" {
				continue
			}
			tagInt, err := strconv.Atoi(tagString)
			if err != nil {
				panic(fmt.Sprintf("marshal error: `surge` tag must be an integer"))
			}
			if _, ok := adt[uint16(tagInt)]; ok {
				panic(fmt.Sprintf("marshal error: `surge` tag must be unique"))
			}
			adt[uint16(tagInt)] = valOf.Field(i).Interface()
		}
		return Marshal(adt, w)
	}

	return newErrUnsupportedMarshalType(v)
}

func Unmarshal(v interface{}, r io.Reader) error {
	switch v := v.(type) {
	case *bool:
		bs := [1]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = bs[0] != 0
		return nil

	case *int8:
		bs := [1]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = int8(bs[0])
		return nil

	case *int16:
		bs := [2]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = int16(binary.LittleEndian.Uint16(bs[:]))
		return nil

	case *int32:
		bs := [4]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = int32(binary.LittleEndian.Uint32(bs[:]))
		return nil

	case *int64:
		bs := [8]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = int64(binary.LittleEndian.Uint64(bs[:]))
		return nil

	case *uint8:
		bs := [1]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = uint8(bs[0])
		return nil

	case *uint16:
		bs := [2]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = binary.LittleEndian.Uint16(bs[:])
		return nil

	case *uint32:
		bs := [4]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = binary.LittleEndian.Uint32(bs[:])
		return nil

	case *uint64:
		bs := [8]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		*v = binary.LittleEndian.Uint64(bs[:])
		return nil

	case *string:
		bs := [4]byte{}
		if _, err := io.ReadFull(r, bs[:]); err != nil {
			return err
		}
		data := make([]byte, binary.LittleEndian.Uint32(bs[:]))
		if _, err := io.ReadFull(r, data); err != nil {
			return err
		}
		*v = string(data)
		return nil
	}

	if i, ok := v.(Unmarshaler); ok {
		return i.Unmarshal(r)
	}

	valOf := reflect.ValueOf(v)
	if valOf.Type().Kind() == reflect.Ptr {
		switch valOf := reflect.Indirect(valOf); valOf.Type().Kind() {
		case reflect.Array, reflect.Slice:
			len := uint32(0)
			if err := Unmarshal(&len, r); err != nil {
				return err
			}
			valOf.Set(reflect.MakeSlice(valOf.Type(), int(len), int(len)))
			for i := 0; i < int(len); i++ {
				if err := Unmarshal(valOf.Index(i).Addr().Interface(), r); err != nil {
					return err
				}
			}
			return nil

		case reflect.Map:
			len := uint32(0)
			if err := Unmarshal(&len, r); err != nil {
				return err
			}
			valOf.Set(reflect.MakeMapWithSize(valOf.Type(), int(len)))
			key := reflect.New(valOf.Type().Key())
			elem := reflect.New(valOf.Type().Elem())
			for i := 0; i < int(len); i++ {
				if err := Unmarshal(key.Interface(), r); err != nil {
					return err
				}
				if err := Unmarshal(elem.Interface(), r); err != nil {
					return err
				}
				valOf.SetMapIndex(reflect.Indirect(key), reflect.Indirect(elem))
			}
			return nil

		case reflect.Struct:
			m := make(map[uint16]reflect.Value, valOf.Type().NumField())
			for i := 0; i < valOf.Type().NumField(); i++ {
				field := valOf.Type().Field(i)
				tag := field.Tag.Get("surge")
				if tag == "" {
					continue
				}
				tagInt, err := strconv.Atoi(tag)
				if err != nil {
					panic(fmt.Sprintf("unmarshal error: `surge` tag must be an integer"))
				}
				if _, ok := m[uint16(tagInt)]; ok {
					panic(fmt.Sprintf("unmarshal error: `surge` tag must be unique"))
				}
				m[uint16(tagInt)] = valOf.Field(i)
			}

			len := uint32(0)
			if err := Unmarshal(&len, r); err != nil {
				return err
			}
			key := uint16(0)
			for i := 0; i < int(len); i++ {
				if err := Unmarshal(&key, r); err != nil {
					return err
				}
				if _, ok := m[key]; !ok {
					return newErrTagNotFound(key)
				}
				elem := reflect.New(m[key].Type())
				if err := Unmarshal(elem.Interface(), r); err != nil {
					return err
				}
				m[key].Set(reflect.Indirect(elem))
			}

			return nil
		}
	}

	return newErrUnsupportedUnmarshalType(v)
}

func SizeHint(v interface{}) int {
	switch v := v.(type) {
	case bool:
		return 1
	case *bool:
		return 1
	case []bool:
		return len(v)

	case int8:
		return 1
	case *int8:
		return 1
	case []int8:
		return len(v)

	case int16:
		return 2
	case *int16:
		return 2
	case []int16:
		return len(v) << 1

	case int32:
		return 4
	case *int32:
		return 4
	case []int32:
		return len(v) << 2

	case int64:
		return 8
	case *int64:
		return 8
	case []int64:
		return len(v) << 3

	case uint8:
		return 1
	case *uint8:
		return 1
	case []uint8:
		return len(v)

	case uint16:
		return 2
	case *uint16:
		return 2
	case []*uint16:
		return len(v) << 1

	case uint32:
		return 4
	case *uint32:
		return 4
	case []*uint32:
		return len(v) << 2

	case uint64:
		return 8
	case *uint64:
		return 8
	case []uint64:
		return len(v) << 3
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
		sizeHint := 4 // Size of length prefix
		for i := 0; i < valOf.Len(); i++ {
			sizeHint += SizeHint(valOf.Index(i).Interface())
		}
		return sizeHint
	}

	panic(newErrUnsupportedSizeHintType(v))
}

func tagStringToTagInt(tag string) uint8 {
	// Finish enumerating this switch statement, because it will be much faster
	// than string parsing.
	switch tag {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "11":
		return 11
	case "12":
		return 12
	case "13":
		return 13
	case "14":
		return 14
	case "15":
		return 15
	case "16":
		return 16
	case "17":
		return 17
	case "18":
		return 18
	case "19":
		return 19
	case "20":
		return 20
	case "21":
		return 21
	case "22":
		return 22
	case "23":
		return 23
	case "24":
		return 24
	case "25":
		return 25
	case "26":
		return 26
	case "27":
		return 27
	case "28":
		return 28
	case "29":
		return 29
	case "30":
		return 30
	case "31":
		return 31
	case "32":
		return 32
	case "33":
		return 33
	case "34":
		return 34
	case "35":
		return 35
	case "36":
		return 36
	case "37":
		return 37
	case "38":
		return 38
	case "39":
		return 39
	case "40":
		return 40
	case "41":
		return 41
	case "42":
		return 42
	case "43":
		return 43
	case "44":
		return 44
	case "45":
		return 45
	case "46":
		return 46
	case "47":
		return 47
	case "48":
		return 48
	case "49":
		return 49
	default:
		panic("exceeded maximum tags: 255")
	}
}
