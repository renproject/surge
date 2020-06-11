package surge_test

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/renproject/surge"
)

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
