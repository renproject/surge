package surge_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// MarshalUnmarshalCheck generates a random instance of a type, marshal it into
// binary, unmarshals the result into a new instance of the type, and then does
// a deep equality check on the two instances. The expectation is that both
// instances will be equal.
func MarshalUnmarshalCheck(t reflect.Type) error {
	// Generate
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal
	data, err := surge.ToBinary(x.Interface())
	if err != nil {
		return fmt.Errorf("cannot marshal: %v", err)
	}
	// Unmarshal
	y := reflect.New(t)
	if err := surge.FromBinary(y.Interface(), data); err != nil {
		return fmt.Errorf("cannot unmarshal: %v", err)
	}
	// Equality
	if !reflect.DeepEqual(x.Interface(), y.Elem().Interface()) {
		return fmt.Errorf("unequal")
	}
	return nil
}

// FuzzCheck generates a random instance of a type, and then attempts to
// unmarshal random bytes into that instance.
func FuzzCheck(t reflect.Type) {
	// Fuzz data
	data, ok := quick.Value(reflect.TypeOf([]byte{}), rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		panic(fmt.Errorf("cannot generate value of type %v", t))
	}
	// Unmarshal
	x := reflect.New(t)
	if err := surge.FromBinary(x.Interface(), data.Bytes()); err != nil {
		// Ignore the error, because we are only interested in whether or not
		// the unmarshaling causes a panic.
	}
}

func MarshalBufTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Generate buffer that is too small
	size := surge.SizeHint(x.Interface())
	// Marshal
	for bufLen := 0; bufLen < size; bufLen++ {
		buf := make([]byte, bufLen)
		rem := size
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

func MarshalRemTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Generate buffer that is too small
	size := surge.SizeHint(x.Interface())
	// Marshal
	for rem := 0; rem < size; rem++ {
		buf := make([]byte, size)
		if _, _, err := surge.Marshal(x.Interface(), buf, rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

func UnmarshalBufTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal the value so that we can attempt to unmarshal the resulting data
	size := surge.SizeHint(x.Interface())
	buf := make([]byte, size)
	rem := size
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when marshaling
		rem = size + 48*x.Len()
	}
	if _, _, err := surge.Marshal(x.Interface(), buf, rem); err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	// Unmarshal with buffers that are too small
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when unmarshaling
		rem = size + x.Len()*int(t.Key().Size()+t.Elem().Size())
	}
	for bufLen := 0; bufLen < size; bufLen++ {
		y := reflect.New(t)
		if _, _, err := surge.Unmarshal(y.Interface(), buf[:bufLen], rem); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

func UnmarshalRemTooSmall(t reflect.Type) error {
	// Generate value
	x, ok := quick.Value(t, rand.New(rand.NewSource(time.Now().UnixNano())))
	if !ok {
		return fmt.Errorf("cannot generate value of type %v", t)
	}
	// Marshal the value so that we can attempt to unmarshal the resulting data
	size := surge.SizeHint(x.Interface())
	buf := make([]byte, size)
	rem := size
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when marshaling
		rem = size + 48*x.Len()
	}
	if _, _, err := surge.Marshal(x.Interface(), buf, rem); err != nil {
		return fmt.Errorf("unexpected error: %v", err)
	}
	// Unmarshal with memory quotas that are too small
	if t.Kind() == reflect.Map {
		// Maps take up extra memory quota when unmarshaling
		rem = size + x.Len()*int(t.Key().Size()+t.Elem().Size())
	}
	for rem2 := 0; rem2 < rem; rem2++ {
		y := reflect.New(t)
		if _, _, err := surge.Unmarshal(y.Interface(), buf, rem2); err == nil {
			return fmt.Errorf("unexpected error: %v", err)
		}
	}
	return nil
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Triangle struct {
	A Point `json:"a"`
	B Point `json:"b"`
	C Point `json:"c"`
}

type Model struct {
	Name      string     `json:"name"`
	Triangles []Triangle `json:"triangles"`
}

func mockPoint() Point {
	return Point{
		X: rand.Float64(),
		Y: rand.Float64(),
		Z: rand.Float64(),
	}
}

func mockTriangle() Triangle {
	return Triangle{
		A: mockPoint(),
		B: mockPoint(),
		C: mockPoint(),
	}
}

func mockModel() Model {
	model := Model{
		Name:      "mock",
		Triangles: make([]Triangle, 100),
	}
	for i := range model.Triangles {
		model.Triangles[i] = mockTriangle()
	}
	return model
}

func BenchmarkPointMarshalJSON(b *testing.B) {
	points := make([]Point, b.N)
	for i := range points {
		points[i] = mockPoint()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&points[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTriangleMarshalJSON(b *testing.B) {
	triangles := make([]Triangle, b.N)
	for i := range triangles {
		triangles[i] = mockTriangle()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&triangles[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkModelMarshalJSON(b *testing.B) {
	models := make([]Model, b.N)
	for i := range models {
		models[i] = mockModel()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(&models[i])
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPointMarshal(b *testing.B) {
	buf := [surge.MaxBytes]byte{}
	points := make([]Point, b.N)
	for i := range points {
		points[i] = mockPoint()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := surge.Marshal(&points[i], buf[:], surge.MaxBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTriangleMarshal(b *testing.B) {
	buf := [surge.MaxBytes]byte{}
	triangles := make([]Triangle, b.N)
	for i := range triangles {
		triangles[i] = mockTriangle()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := surge.Marshal(&triangles[i], buf[:], surge.MaxBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkModelMarshal(b *testing.B) {
	buf := [surge.MaxBytes]byte{}
	models := make([]Model, b.N)
	for i := range models {
		models[i] = mockModel()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := surge.Marshal(&models[i], buf[:], surge.MaxBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type Foo struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func (foo Foo) Marshal(buf []byte, rem int) ([]byte, int, error) {
	var err error
	if buf, rem, err = surge.MarshalString(foo.Name, buf, rem); err != nil {
		return buf, rem, err
	}
	if buf, rem, err = surge.MarshalI64(foo.BirthDay.UnixNano(), buf, rem); err != nil {
		return buf, rem, err
	}
	if buf, rem, err = surge.MarshalString(foo.Phone, buf, rem); err != nil {
		return buf, rem, err
	}
	if buf, rem, err = surge.MarshalI64(int64(foo.Siblings), buf, rem); err != nil {
		return buf, rem, err
	}
	if buf, rem, err = surge.MarshalBool(foo.Spouse, buf, rem); err != nil {
		return buf, rem, err
	}
	if buf, rem, err = surge.MarshalF64(foo.Money, buf, rem); err != nil {
		return buf, rem, err
	}
	return buf, rem, err
}

func BenchmarkFoo(b *testing.B) {
	buf := [surge.MaxBytes]byte{}
	foos := make([]Foo, b.N)
	for i := range foos {
		foos[i] = Foo{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := foos[i].Marshal(buf[:], surge.MaxBytes)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type Bar int64

func (Bar) SizeHint() int {
	return 42
}

func (Bar) Marshal(buf []byte, rem int) ([]byte, int, error) {
	copy(buf, []byte{0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42, 0x42})
	return buf[:42], rem - 42, nil
}

func (bar *Bar) Unmarshal(buf []byte, rem int) ([]byte, int, error) {
	bytes := [42]byte{}
	copy(bytes[:], buf[:42])
	*bar = 42
	return buf[:42], rem - 42, nil
}

var _ = Describe("Size hint", func() {
	Context("when a custom implementation exists", func() {
		It("should use the custom implementation", func() {
			Expect(surge.SizeHint(Bar(42))).To(Equal(Bar(42).SizeHint()))
		})
	})
})

var _ = Describe("Marshal", func() {
	Context("when a custom implementation exists", func() {
		It("should use the custom implementation", func() {
			buf := [42]byte{}
			rem := 42
			_, _, err := surge.Marshal(Bar(42), buf[:], rem)
			Expect(err).ToNot(HaveOccurred())

			buf2 := [42]byte{}
			rem2 := 42
			_, _, err = Bar(42).Marshal(buf2[:], rem2)
			Expect(err).ToNot(HaveOccurred())

			Expect(bytes.Equal(buf[:], buf2[:])).To(BeTrue())
		})
	})

	Context("when marshaling pointers", func() {
		It("should be the same as marshaling the underlying type", func() {
			f := func(x string) bool {
				buf := make([]byte, surge.SizeHint(x))
				rem := surge.SizeHint(x)
				_, _, err := surge.Marshal(x, buf, rem)
				Expect(err).ToNot(HaveOccurred())

				buf2 := make([]byte, surge.SizeHint(&x))
				rem2 := surge.SizeHint(&x)
				_, _, err = surge.Marshal(&x, buf2, rem2)
				Expect(err).ToNot(HaveOccurred())

				Expect(bytes.Equal(buf, buf2)).To(BeTrue())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})

var _ = Describe("Unmarshal", func() {
	Context("when a custom implementation exists", func() {
		It("should use the custom implementation", func() {
			bar := Bar(0)
			buf := [42]byte{}
			rem := 42
			_, _, err := surge.Unmarshal(&bar, buf[:], rem)
			Expect(err).ToNot(HaveOccurred())

			bar2 := Bar(1)
			buf2 := [42]byte{}
			rem2 := 42
			_, _, err = bar2.Unmarshal(buf2[:], rem2)
			Expect(err).ToNot(HaveOccurred())

			Expect(bar).To(Equal(bar2))
		})
	})

	Context("when unmarshaling non-pointers", func() {
		It("should return an error", func() {
			f := func(x string) bool {
				buf := make([]byte, surge.SizeHint(x))
				rem := surge.SizeHint(x)
				_, _, err := surge.Unmarshal(x, buf, rem)
				Expect(err).To(HaveOccurred())
				return true
			}
			Expect(quick.Check(f, nil)).To(Succeed())
		})
	})
})
