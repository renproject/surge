# `🔌 surge`

[![GitHub](https://github.com/renproject/surge/workflows/test/badge.svg)](https://github.com/renproject/surge/workflows/test/badge.svg)
[![Coverage](https://coveralls.io/repos/github/renproject/surge/badge.svg?branch=master)](https://coveralls.io/github/renproject/surge?branch=master)
[![Report](https://goreportcard.com/badge/github.com/renproject/surge)](https://goreportcard.com/badge/github.com/renproject/surge)

[Documentation](https://godoc.org/github.com/renproject/surge)

A library for fast binary (un)marshaling. Designed to be used in Byzantine networks, `🔌 surge` never explicitly panics, protects against malicious inputs, allocates minimally, and has very few dependencies (its only dependency is the [`ginkgo`](https://onsi.github.io/ginkgo) testing framework). It supports the (un)marshaling of:

- [x] scalars,
- [x] arrays,
- [x] slices,
- [x] maps, and
- [x] custom implementations (using the `Marshaler` and `Unmarshaler` interfaces).

## Built-in Types

All built-in types that can be marshaled are supported by `surge`. And, for the vast majority of use cases, `ToBinary` and `FromBinary` are the only functions that you will need to use:

### Scalars

```go
// Marshal
x := uint64(42)
data, err := surge.ToBinary(x)
if err != nil {
    panic(err)
}

// Unmarshal
y := uint64(0)
if err := surge.FromBinary(&y, data); err != nil {
    panic(err)
}
```

### Arrays

Arrays are collections of a known, fixed, length. Arrays are *not* length prefixed, because their length is part of their type. All arrays marshal their elements one-by-one, with the exception of byte arrays (which are marshaled in bulk using `copy`):

```go
// Marshal
x := [4]uint64{42, 43, 44, 45}
data, err := surge.ToBinary(x)
if err != nil {
    panic(err)
}

// Unmarshal
y := [4]uint64{}
if err := surge.FromBinary(&y, data); err != nil {
    panic(err)
}
```

### Slices

Slices are collections of variable length. Slices are length prefixed, because their length is not known at compile-time. All slices marshal their elements one-by-one, with the exception of byte slices (which are marshaled in bulk using `copy`):

```go
// Marshal
x := []uint64{42, 43, 44, 45}
data, err := surge.ToBinary(x)
if err != nil {
    panic(err)
}

// Unmarshal
y := []uint64{}
if err := surge.FromBinary(&y, data); err != nil {
    panic(err)
}
```

### Maps

Maps are effectively slices of key/value pairs. Maps are length prefixed, because their length is not known at compile-time. Maps are marshaled as a sorted slice of `(key, value)` tuples, sorted lexographically by keys (after the key has been marshaled, because not all key types are directly comparable). Sorting is done because it guarantees that the binary output is always the same when the key/value pairs are the same (this is particularly useful when hashing/signing maps for authenticity):

```go
// Marshal
x := map[string]uint64{"foo": 42, "bar": 43, "baz": 44}
data, err := surge.ToBinary(x)
if err != nil {
    panic(err)
}

// Unmarshal
y := map[string]uint64{}
if err := surge.FromBinary(&y, data); err != nil {
    panic(err)
}
```

## User-defined types

The same pattern that we have seen above works for custom structs too. You will *not* need to make any changes to your struct, as long as all of its fields are marshalable by `surge`:

```go
type MyStruct struct {
    Foo int64
    Bar float64
    Baz MyInnerStruct
}

type MyInnerStruct struct {
    Inner1 []bool
    Inner2 []string
}

// Marshal
x := MyStruct{
    Foo: int64(43),
    Bar: float64(3.14),
    Baz: MyInnerStruct{
        Inner1: []bool{true, false},
        Inner2: []string{"hello", "world"},
    },
}
data, err := surge.ToBinary(x)
if err != nil {
    panic(err)
}

// Unmarshal
y := MyStruct{}
if err := surge.FromBinary(&y, data); err != nil {
    panic(err)
}
```

### Specialisation

Using the default marshaler built into `surge` is great for prototyping, and will good enough for many applications. But, sometimes we need to specialise our marshaling. Providing our own implementation will not only be faster, but it will also give us the ability to customise the marshaler (which can be necessary when thinking about backward compatibility, etc.):

```go
type MyStruct struct {
  Foo int64
  Bar float64
  Baz string
}

// SizeHint tells surge how many bytes our
// custom type needs when being represented
// in its binary form.
func (myStruct MyStruct) SizeHint() int {
    return surge.SizeHintI64 +
           surge.SizeHintF64 +
           surge.SizeHintString(myStruct.Baz)
}

// Marshal tells surge exactly how to marshal
// our custom type. As you can see, most implementations
// will be very straight forward, and mostly exist
// for performance reasons. In the future, surge might
// adopt some kind of generator to automatically
// generate these implementations.
func (myStruct MyStruct) Marshal(buf []byte, rem int) ([]byte, int, error) {
    var err error
    if buf, rem, err = surge.MarshalI64(myStruct.Foo, buf, rem); err != nil {
        return buf, rem, err
    }
    if buf, rem, err = surge.MarshalF64(myStruct.Bar, buf, rem); err != nil {
        return buf, rem, err
    }
    if buf, rem, err = surge.MarshalString(myStruct.Baz, buf, rem); err != nil {
        return buf, rem, err
    }
    return buf, rem, err
}

// Unmarshal is the opposite of Marshal, and requires
// a pointer receiver.
func (myStruct *MyStruct) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
    var err error
    if buf, rem, err = surge.UnmarshalI64(&myStruct.Foo, buf, rem); err != nil {
        return buf, rem, err
    }
    if buf, rem, err = surge.UnmarshalF64(&myStruct.Bar, buf, rem); err != nil {
        return buf, rem, err
    }
    if buf, rem, err = surge.UnmarshalString(&myStruct.Baz, buf, rem); err != nil {
        return buf, rem, err
    }
    return buf, rem, err
}
```

### Testing

Testing custom marshaling implementations is incredibly important, but it can also be very tedious, and so it is rarely done as extensively as it should be. Luckily, `surge` helps us get this done quickly. By using the `surgeutil` package, we can write comprehensive tests very quickly:

```go
func TestMyStruct(t *testing.T) {
    // Reflect on our custom type
    t := reflect.TypeOf(MyStruct{})
    
    // Fuzz and expect that it does not panic.
    surgeutil.Fuzz(t)
    
    // Marshal, then unmarshal, then check for
    // equality, and expect there to be no
    // errors.
    if err := surgeutil.MarshalUnmarshalCheck(t); err != nil {
        t.Fatalf("bad marshal/unmarshal/check: %v", err)
    }
    
    // Marshal when the buffer is too small
    // and check that it does not work.
    if err := surgeutil.MarshalBufTooSmall(t); err != nil {
        t.Fatalf("bad marshal with insufficient buffer: %v", err)
    }
    
    // Marshal when the remaining memory quota
    // is too small and check that it does not
    // work.
    if err := surgeutil.MarshalRemTooSmall(t); err != nil {
        t.Fatalf("bad marshal with insufficient rem quota: %v", err)
    }
    
    // Unmarshal when the buffer is too small
    // and check that it does not work.
    if err := surgeutil.UnmarshalBufTooSmall(t); err != nil {
        t.Fatalf("bad marshal with insufficient buffer: %v", err)
    }
    
    // Unmarshal when the remaining memory quota
    // is too small and check that it does not
    // work.
    if err := surgeutil.UnmarshalRemTooSmall(t); err != nil {
        t.Fatalf("bad marshal with insufficient rem quota: %v", err)
    }
}
```

Internally, `surgeutil` makes use of the [`quick`](https://golang.org/pkg/testing/quick) standard library. So, for `surgeutil` to work, your type needs to be compatible with `quick`. This is usually automatic, and most of the time you will not need to think about `quick` at all. For the more exotic types, that do need custom support, all you need to do is implement the [`quick.Generator`](https://golang.org/pkg/testing/quick/#Generator) interface. For more examples of `surgeutil` in use, checkout any of the `*_test.go` files. All of the testing in `surge` is done using the `surgeutil` package.

## Benchmarks

When using specialised implementations, `surge` is about as fast as you can get; it does not really do much under-the-hood. When using the default implementations, the need to use `reflect` introduces some slow-down, but performance is still faster than most alternatives:

```
goos: darwin
goarch: amd64
pkg: github.com/renproject/surge
BenchmarkPointMarshalJSON-8              2064483               563 ns/op              80 B/op          1 allocs/op
BenchmarkTriangleMarshalJSON-8            583173              1752 ns/op             239 B/op          1 allocs/op
BenchmarkModelMarshalJSON-8                 7018            163255 ns/op           24588 B/op          1 allocs/op
BenchmarkPointMarshal-8                 11212546               109 ns/op               0 B/op          0 allocs/op
BenchmarkTriangleMarshal-8               3700579               294 ns/op               0 B/op          0 allocs/op
BenchmarkModelMarshal-8                    38652             28270 ns/op              32 B/op          1 allocs/op
BenchmarkFoo-8                          33130609                33 ns/op               0 B/op          0 allocs/op
```

## Contributions

Built with ❤ by Ren. 
