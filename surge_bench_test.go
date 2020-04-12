package surge_test

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/renproject/surge"
)

func BenchmarkBool(b *testing.B) {
	y := bool(false)
	bin, _ := ToBinary(true)
	for n := 0; n < b.N; n++ {
		bin, _ = ToBinary(true)
		FromBinary(bin, &y)
	}
}

func BenchmarkInt64(b *testing.B) {
	y := int64(1)
	bin, _ := ToBinary(10)
	for n := 0; n < b.N; n++ {
		bin, _ = ToBinary(10)
		FromBinary(bin, &y)
	}
}

func BenchmarkByteSlice(b *testing.B) {
	xs := make([]byte, 100)
	for i := uint16(0); i < 100; i++ {
		xs[i] = byte(rand.Int63())
	}
	ys := make([]byte, 0)

	bin, _ := ToBinary(xs)
	for n := 0; n < b.N; n++ {
		bin, _ = ToBinary(xs)
		FromBinary(bin, &ys)
	}
}

func BenchmarkInt64Slice(b *testing.B) {
	xs := make([]int64, 100)
	for i := uint16(0); i < 100; i++ {
		xs[i] = int64(rand.Int63())
	}
	ys := make([]int64, 0)

	bin, _ := ToBinary(xs)
	for n := 0; n < b.N; n++ {
		bin, _ = ToBinary(xs)
		FromBinary(bin, &ys)
	}
}

func BenchmarkByteArray(b *testing.B) {
	x := [3]byte{byte(rand.Int63()), byte(rand.Int63()), byte(rand.Int63())}
	y := [3]byte{}
	bin, _ := ToBinary(x)
	for n := 0; n < b.N; n++ {
		bin, _ = ToBinary(x)
		FromBinary(bin, &y)
	}
}

func BenchmarkCustomStruct(b *testing.B) {
	a := Point{1, 2}
	x := Point{}
	for n := 0; n < b.N; n++ {
		bin, _ := ToBinary(a)
		FromBinary(bin, &x)
	}
}

func BenchmarkMapStringUint64(b *testing.B) {
	xs := make(map[string]uint64)
	for i := uint16(0); i < 100; i++ {

		xs[string(fmt.Sprintf("%d", rand.Uint64()))] = rand.Uint64()
	}
	ys := make(map[string]uint64)

	for n := 0; n < b.N; n++ {
		bin, _ := ToBinary(&xs)
		FromBinary(bin, &ys)
	}
}
