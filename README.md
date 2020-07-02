# `🔌 surge`

[![GitHub](https://github.com/renproject/surge/workflows/test/badge.svg)](https://github.com/renproject/surge/workflows/test/badge.svg)
[![Coverage](https://coveralls.io/repos/github/renproject/surge/badge.svg?branch=master)](https://coveralls.io/github/renproject/surge?branch=master)
[![Report](https://goreportcard.com/badge/github.com/renproject/surge)](https://goreportcard.com/badge/github.com/renproject/surge)

[Documentation](https://godoc.org/github.com/renproject/surge)

A library for fast binary (un)marshaling. Designed to be used in Byzantine networks, `🔌 surge` never explicitly panics and protects against malicious inputs. It supports the (un)marshaling of:

- [x] scalars,
- [x] arrays,
- [x] slices,
- [x] maps, and
- [x] custom implementations (using the `Marshaler` and `Unmarshaler` interfaces).

## Examples

For most use cases, `ToBinary` and `FromBinary` are the only functions that you will need to interact with:

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

The same pattern works for custom struct too, without any changes:

```go
type MyStruct struct {
  Foo int64
  Bar float64
  Baz string
}

// Marshal
x := MyStruct{ Foo: int64(43), Bar: float64(3.14), Baz: "baz" }
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

## Custom Implementation

Using the default marshaler built into `surge` is great for prototyping, and will good enough for many applications. But, sometimes we need to specialise our marshaling. Providing our own implementation will not only be faster, but it will also give us the ability to customise the marshaler (which can be necessary when thinking about backward compatibility, etc.):

```
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
BenchmarkFoo-8                          33130609                33.0 ns/op             0 B/op          0 allocs/op
```

## Contributions

Built with ❤ by Ren. 
