package surge_test

import (
	"bytes"
	"io"
	"math/rand"
	"reflect"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/surge"
)

type Point struct {
	x uint64
	y uint64
}

func (p Point) SizeHint() int {
	return SizeHint(p.x) + SizeHint(p.y)
}

func (p Point) Marshal(w io.Writer, m int) (int, error) {
	m, err := Marshal(w, p.x, m)
	if err != nil {
		return m, err
	}
	return Marshal(w, p.y, m)
}

func (p *Point) Unmarshal(r io.Reader, m int) (int, error) {
	m, err := Unmarshal(r, &p.x, m)
	if err != nil {
		return m, err
	}
	return Unmarshal(r, &p.y, m)
}

var _ = Describe("Surge", func() {

	Context("when marshaling bool", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x bool) bool {
					y := bool(false)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling int8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x int8) bool {
					y := int8(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling int16", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x int16) bool {
					y := int16(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling int32", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x int32) bool {
					y := int32(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling int64", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x int64) bool {
					y := int64(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling uint8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x uint8) bool {
					y := uint8(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling uint16", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x uint16) bool {
					y := uint16(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling uint32", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x uint32) bool {
					y := uint32(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling uint64", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x uint64) bool {
					y := uint64(0)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling []byte", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					xs := make([]byte, n)
					for i := uint16(0); i < n; i++ {
						xs[i] = byte(rand.Int63())
					}
					ys := make([]byte, 0)
					bin, err := ToBinary(xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &ys)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling string", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					data := make([]byte, n)
					for i := uint16(0); i < n; i++ {
						data[i] = byte(rand.Int63())
					}
					xs := string(data)
					ys := string("")
					bin, err := ToBinary(xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &ys)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling *uint64", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x uint64) bool {
					y := uint64(0)
					bin, err := ToBinary(&x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &y)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling []uint64", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					xs := make([]uint64, n)
					for i := uint16(0); i < n; i++ {
						xs[i] = rand.Uint64()
					}
					ys := make([]uint64, 0)
					bin, err := ToBinary(xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &ys)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling []uint64 array", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				emptyxs := []uint64{}
				emptyys := []uint64{}
				bin, err := ToBinary(emptyxs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(bin, &emptyys)).ToNot(HaveOccurred())
				Expect(emptyxs).To(Equal(emptyys))

				xs := [1000]uint64{}
				for i := 0; i < 1000; i++ {
					xs[i] = rand.Uint64()
				}
				ys := [1000]uint64{}
				bin, err = ToBinary(xs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(bin, &ys)).ToNot(HaveOccurred())
				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling map[uint64]uint64", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					xs := make(map[uint64]uint64)
					for i := uint16(0); i < n; i++ {
						xs[rand.Uint64()] = rand.Uint64()
					}
					ys := make(map[uint64]uint64)
					bin, err := ToBinary(&xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &ys)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling a struct", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x, y uint64) bool {
					a := Point{x, y}
					b := Point{}
					bin, err := ToBinary(a)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(bin, &b)).ToNot(HaveOccurred())
					return reflect.DeepEqual(a, b)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling with positive capacity", func() {
		Context("when marshaling too many bytes", func() {
			Context("when marshaling []byte", func() {
				It("should return a negative capacity", func() {
					f := func(n uint16) bool {
						xs := make([]byte, n+2)
						for i := 0; i < int(n)+2; i++ {
							xs[i] = byte(rand.Int63())
						}
						m, err := Marshal(new(bytes.Buffer), xs, 1)
						Expect(m).To(BeNumerically("<", 0))
						Expect(err).ToNot(HaveOccurred())
						m, err = Marshal(new(bytes.Buffer), xs, int(n)+1)
						Expect(m).To(BeNumerically("<", 0))
						Expect(err).ToNot(HaveOccurred())
						return true
					}
					err := quick.Check(f, nil)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when marshaling string", func() {
				It("should return a negative capacity", func() {
					f := func(n uint16) bool {
						data := make([]byte, int(n)+2)
						for i := 0; i < int(n)+2; i++ {
							data[i] = byte(rand.Int63())
						}
						xs := string(data)
						m, err := Marshal(new(bytes.Buffer), xs, 1)
						Expect(m).To(BeNumerically("<", 0))
						Expect(err).ToNot(HaveOccurred())
						m, err = Marshal(new(bytes.Buffer), xs, int(n)+1)
						Expect(m).To(BeNumerically("<", 0))
						Expect(err).ToNot(HaveOccurred())
						return true
					}
					err := quick.Check(f, nil)
					Expect(err).ToNot(HaveOccurred())
				})
			})

			Context("when marshaling a struct", func() {
				It("should return an error", func() {
					f := func(x, y uint64) bool {
						a := Point{x, y}
						_, err := Marshal(new(bytes.Buffer), a, 1)
						Expect(err).To(HaveOccurred())
						_, err = Marshal(new(bytes.Buffer), a, 7)
						Expect(err).To(HaveOccurred())
						return true
					}
					err := quick.Check(f, nil)
					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})

	Context("when marshaling with negative capacity", func() {
		It("should return an error", func() {
			m, err := Marshal(nil, nil, 0)
			Expect(m).To(Equal(0))
			Expect(err).To(HaveOccurred())
			m, err = Marshal(nil, nil, -1)
			Expect(m).To(Equal(-1))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when marshaling an unsupported type", func() {
		It("should return an error", func() {
			_, err := ToBinary(struct{}{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when unmarshaling with negative capacity", func() {
		It("should return an error", func() {
			m, err := Unmarshal(nil, nil, 0)
			Expect(m).To(Equal(0))
			Expect(err).To(HaveOccurred())
			m, err = Unmarshal(nil, nil, -1)
			Expect(m).To(Equal(-1))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when marshaling an unsupported type", func() {
		It("should return an error", func() {
			err := FromBinary([]byte{}, &struct{}{})
			Expect(err).To(HaveOccurred())
		})
	})

	Context("when unmarshaling to a non-pointer type", func() {
		It("should return an error", func() {
			y := uint64(0)
			err := FromBinary([]byte{14, 26, 37, 48}, y)
			Expect(err).To(HaveOccurred())
		})
	})
})
