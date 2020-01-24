package surge

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
	"sort"
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

// ToBinary marshals a value into a byte slice. It allocates an in-memory
// buffer, using `SizeHint` to estimate the initial size of the buffer
// (preventing the need to grow the buffer during marshaling). This function
// should not be used to implement the `Marshaler` interface, or as part of
// marshaling a parent value.
//
// An example of marshaling a scalar value:
//
//      x := 42
//      data, err := surge.ToBinary(x)
//      if err != nil {
//          panic(err)
//      }
//      fmt.Printf("%x", data)
//
// An example of marshaling a custom struct:
//
//		type Point struct {
//			X uint64 `surge:"0"`
//			Y uint64 `surge:"1"`
//		}
//
//		p := Point{ 13, 169 }
//		data, err := surge.ToBinary(p)
//      if err != nil {
//          panic(err)
//      }
//      fmt.Printf("%x", data)
//
func ToBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(SizeHint(v))
	err := Marshal(v, buf)
	return buf.Bytes(), err
}

// FromBinary unmarshals a byte slice into a destination value. The destination
// value must be a pointer.
//
// An example of marshaling/unmarshaling a map:
//
//      xs := map[string]string{}
//		xs["foo1"] = "bar"
//		xs["foo2"] = "baz"
//
//		data, err := surge.ToBinary(xs)
//		if err != nil {
//			panic(err)
//		}
//
//      ys := map[string]string{}
//      if err := surge.FromBinary(&ys, data); err != nil {
//          panic(err)
//      }
//		fmt.Printf("foo1: %s\n", ys["foo1"])
//		fmt.Printf("foo2: %s\n", ys["foo2"])
//
func FromBinary(v interface{}, data []byte) error {
	buf := bytes.NewBuffer(data)
	return Unmarshal(v, buf)
}

// Marshal a value into an `io.Writer`. This is more efficient than marshaling a
// value into a byte slice and returning the byte slice, because the `io.Writer`
// can be pre-allocated with enough memory to avoid internal allocations and
// buffer copies. This function should only used when defining custom
// implementations of the `Marshaler` interface, or as part of marshaling a
// parent value. In most use cases, the `ToBinary` function should be used.
//
// - When marshaling scalars, all values are marshaled into bytes using little
//   endian encoding.
// - When marshaling arrays/slices/maps, an uint32 length prefix is marshaled
//   and prefixed.
// - When marshaling maps, key/value pairs are marshaled in order of the keys
//   (sorted after the key has been marshaled). This guarantees consistency; the
//   marshaled bytes are always the same if the key/values in the map are the
//   same. This is particularly useful when hashing.
// - When marshaling custom struct, the ``surge:`` struct tags are used to
//   convert the struct into a map (which is then marshaled like a normal map).
// - When marshaling a value that implements the `Marshaler` interface, it is up
//   to the user to guarantee that the implementation is sane.
func Marshal(v interface{}, w io.Writer) error {
	// Marshal scalar types.
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
		adt := make(map[uint8]interface{}, valOf.Type().NumField())
		for i := 0; i < valOf.Type().NumField(); i++ {
			field := valOf.Type().Field(i)
			tagString := field.Tag.Get("surge")
			if tagString == "" {
				continue
			}
			tagInt := FastStrconvUint8(tagString)
			if tagInt == 255 {
				return newErrTagOutOfRange(tagInt)
			}
			if _, ok := adt[tagInt]; ok {
				return newErrTagDuplicate(tagInt)
			}
			adt[tagInt] = valOf.Field(i).Interface()
		}
		return Marshal(adt, w)
	}

	return newErrUnsupportedMarshalType(v)
}

// Unmarshal from an `io.Reader` into a destination value. The destination value
// must be a pointer. This function should only used when defining custom
// implementations of the `Unmarshaler` interface, or as part of marshaling a
// parent value. In most use cases, the `FromBinary` function should be used.
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
		case reflect.Array:
			len := uint32(0)
			if err := Unmarshal(&len, r); err != nil {
				return err
			}
			if uint32(valOf.Len()) != len {
				return newErrBadLength(uint32(valOf.Len()), len)
			}
			for i := 0; i < int(len); i++ {
				if err := Unmarshal(valOf.Index(i).Addr().Interface(), r); err != nil {
					return err
				}
			}
			return nil

		case reflect.Slice:
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
			m := make(map[uint8]reflect.Value, valOf.Type().NumField())
			for i := 0; i < valOf.Type().NumField(); i++ {
				field := valOf.Type().Field(i)
				tagString := field.Tag.Get("surge")
				if tagString == "" {
					continue
				}
				tagInt := FastStrconvUint8(tagString)
				if tagInt == 255 {
					return newErrTagOutOfRange(tagInt)
				}
				if _, ok := m[tagInt]; ok {
					return newErrTagDuplicate(tagInt)
				}
				m[tagInt] = valOf.Field(i)
			}

			len := uint32(0)
			if err := Unmarshal(&len, r); err != nil {
				return err
			}
			key := uint8(0)
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
		return 4 + len(v)<<1

	case int32:
		return 4
	case *int32:
		return 4
	case []int32:
		return 4 + len(v)<<2

	case int64:
		return 8
	case *int64:
		return 8
	case []int64:
		return 4 + len(v)<<3

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
		return 4 + len(v)<<1

	case uint32:
		return 4
	case *uint32:
		return 4
	case []*uint32:
		return 4 + len(v)<<2

	case uint64:
		return 8
	case *uint64:
		return 8
	case []uint64:
		return 4 + len(v)<<3
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

	case reflect.Map:
		sizeHint := 4 // Size of length prefix
		for _, key := range valOf.MapKeys() {
			sizeHint += SizeHint(key.Interface())
			sizeHint += SizeHint(valOf.MapIndex(key).Interface())
		}
		return sizeHint

	case reflect.Struct:
		adt := make(map[uint8]struct{}, valOf.Type().NumField())
		sizeHint := 4 // Size of length prefix
		for i := 0; i < valOf.Type().NumField(); i++ {
			field := valOf.Type().Field(i)
			tagString := field.Tag.Get("surge")
			if tagString == "" {
				continue
			}
			tagInt := FastStrconvUint8(tagString)
			if tagInt == 255 {
				// Ignore size hinting for malformed tags.
				continue
			}
			if _, ok := adt[tagInt]; ok {
				// Ignore size hinting for non-unique tags.
				continue
			}
			adt[tagInt] = struct{}{}
			sizeHint += 1
			sizeHint += SizeHint(valOf.Field(i).Interface())
		}
		return SizeHint(adt)
	}

	return 0
}

// FastStrconvUint8 uses a large switch statement to speed up the conversion of
// strings to uint8s for the specific use-case of tags. A return value of 255
// should be interpreted as an error.
func FastStrconvUint8(tag string) uint8 {
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
	case "50":
		return 50
	case "51":
		return 51
	case "52":
		return 52
	case "53":
		return 53
	case "54":
		return 54
	case "55":
		return 55
	case "56":
		return 56
	case "57":
		return 57
	case "58":
		return 58
	case "59":
		return 59
	case "60":
		return 60
	case "61":
		return 61
	case "62":
		return 62
	case "63":
		return 63
	case "64":
		return 64

	case "65":
		return 65
	case "66":
		return 66
	case "67":
		return 67
	case "68":
		return 68
	case "69":
		return 69
	case "70":
		return 70
	case "71":
		return 71
	case "72":
		return 72
	case "73":
		return 73
	case "74":
		return 74
	case "75":
		return 75
	case "76":
		return 76
	case "77":
		return 77
	case "78":
		return 78
	case "79":
		return 79
	case "80":
		return 80
	case "81":
		return 81
	case "82":
		return 82
	case "83":
		return 83
	case "84":
		return 84
	case "85":
		return 85
	case "86":
		return 86
	case "87":
		return 87
	case "88":
		return 88
	case "89":
		return 89
	case "90":
		return 90
	case "91":
		return 91
	case "92":
		return 92
	case "93":
		return 93
	case "94":
		return 94
	case "95":
		return 95
	case "96":
		return 96
	case "97":
		return 97
	case "98":
		return 98
	case "99":
		return 99
	case "100":
		return 100
	case "101":
		return 101
	case "102":
		return 102
	case "103":
		return 103
	case "104":
		return 104
	case "105":
		return 105
	case "106":
		return 106
	case "107":
		return 107
	case "108":
		return 108
	case "109":
		return 109
	case "110":
		return 110
	case "111":
		return 111
	case "112":
		return 112
	case "113":
		return 113
	case "114":
		return 114
	case "115":
		return 115
	case "116":
		return 116
	case "117":
		return 117
	case "118":
		return 118
	case "119":
		return 119
	case "120":
		return 120
	case "121":
		return 121
	case "122":
		return 122
	case "123":
		return 123
	case "124":
		return 124
	case "125":
		return 125
	case "126":
		return 126
	case "127":
		return 127
	case "128":
		return 128

	case "129":
		return 129
	case "130":
		return 130
	case "131":
		return 131
	case "132":
		return 132
	case "133":
		return 133
	case "134":
		return 134
	case "135":
		return 135
	case "136":
		return 136
	case "137":
		return 137
	case "138":
		return 138
	case "139":
		return 139
	case "140":
		return 140
	case "141":
		return 141
	case "142":
		return 142
	case "143":
		return 143
	case "144":
		return 144
	case "145":
		return 145
	case "146":
		return 146
	case "147":
		return 147
	case "148":
		return 148
	case "149":
		return 149
	case "150":
		return 150
	case "151":
		return 151
	case "152":
		return 152
	case "153":
		return 153
	case "154":
		return 154
	case "155":
		return 155
	case "156":
		return 156
	case "157":
		return 157
	case "158":
		return 158
	case "159":
		return 159
	case "160":
		return 160
	case "161":
		return 161
	case "162":
		return 162
	case "163":
		return 163
	case "164":
		return 164
	case "165":
		return 165
	case "166":
		return 166
	case "167":
		return 167
	case "168":
		return 168
	case "169":
		return 169
	case "170":
		return 170
	case "171":
		return 171
	case "172":
		return 172
	case "173":
		return 173
	case "174":
		return 174
	case "175":
		return 175
	case "176":
		return 176
	case "177":
		return 177
	case "178":
		return 178
	case "179":
		return 179
	case "180":
		return 180
	case "181":
		return 181
	case "182":
		return 182
	case "183":
		return 183
	case "184":
		return 184
	case "185":
		return 185
	case "186":
		return 186
	case "187":
		return 187
	case "188":
		return 188
	case "189":
		return 189
	case "190":
		return 190
	case "191":
		return 191
	case "192":
		return 192
	case "193":
		return 193

	case "194":
		return 194
	case "195":
		return 195
	case "196":
		return 196
	case "197":
		return 197
	case "198":
		return 198
	case "199":
		return 199
	case "200":
		return 200
	case "201":
		return 201
	case "202":
		return 202
	case "203":
		return 203
	case "204":
		return 204
	case "205":
		return 205
	case "206":
		return 206
	case "207":
		return 207
	case "208":
		return 208
	case "209":
		return 209
	case "210":
		return 210
	case "211":
		return 211
	case "212":
		return 212
	case "213":
		return 213
	case "214":
		return 214
	case "215":
		return 215
	case "216":
		return 216
	case "217":
		return 217
	case "218":
		return 218
	case "219":
		return 219
	case "220":
		return 220
	case "221":
		return 221
	case "222":
		return 222
	case "223":
		return 223
	case "224":
		return 224
	case "225":
		return 225
	case "226":
		return 226
	case "227":
		return 227
	case "228":
		return 228
	case "229":
		return 229
	case "230":
		return 230
	case "231":
		return 231
	case "232":
		return 232
	case "233":
		return 233
	case "234":
		return 234
	case "235":
		return 235
	case "236":
		return 236
	case "237":
		return 237
	case "238":
		return 238
	case "239":
		return 239
	case "240":
		return 240
	case "241":
		return 241
	case "242":
		return 242
	case "243":
		return 243
	case "244":
		return 244
	case "245":
		return 245
	case "246":
		return 246
	case "247":
		return 247
	case "248":
		return 248
	case "249":
		return 249
	case "250":
		return 250
	case "251":
		return 251
	case "252":
		return 252
	case "253":
		return 253
	case "254":
		return 254

	default:
		return 255
	}
}
