package surge_test

import (
	"math/rand"
	"reflect"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/surge"
)

var _ = Describe("Surge", func() {
	Context("when marshaling bool", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x bool) bool {
					y := bool(false)
					bin, err := ToBinary(x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
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
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling *int8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(x int8) bool {
					y := int8(0)
					bin, err := ToBinary(&x)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(&y, bin)).ToNot(HaveOccurred())
					return reflect.DeepEqual(x, y)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling []int8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					xs := make([]int8, n)
					for i := uint16(0); i < n; i++ {
						xs[i] = int8(rand.Int63())
					}
					ys := make([]int8, 0)
					bin, err := ToBinary(xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(&ys, bin)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling *[]int8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint16) bool {
					xs := make([]int8, n)
					for i := uint16(0); i < n; i++ {
						xs[i] = int8(rand.Int63())
					}
					ys := make([]int8, 0)
					bin, err := ToBinary(&xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(&ys, bin)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("when marshaling []int8 array", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				emptyxs := []int8{}
				emptyys := []int8{}
				bin, err := ToBinary(emptyxs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(&emptyys, bin)).ToNot(HaveOccurred())
				Expect(emptyxs).To(Equal(emptyys))

				xs := [1000]int8{}
				for i := 0; i < 1000; i++ {
					xs[i] = int8(rand.Int63())
				}
				ys := [1000]int8{}
				bin, err = ToBinary(xs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(&ys, bin)).ToNot(HaveOccurred())
				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling *[]int8 array", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				emptyxs := []int8{}
				emptyys := []int8{}
				bin, err := ToBinary(&emptyxs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(&emptyys, bin)).ToNot(HaveOccurred())
				Expect(emptyxs).To(Equal(emptyys))

				xs := [1000]int8{}
				for i := 0; i < 1000; i++ {
					xs[i] = int8(rand.Int63())
				}
				ys := [1000]int8{}
				bin, err = ToBinary(&xs)
				Expect(err).ToNot(HaveOccurred())
				Expect(FromBinary(&ys, bin)).ToNot(HaveOccurred())
				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling map[int8]int8", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				f := func(n uint8) bool {
					xs := make(map[int8]int8)
					for i := uint8(0); i < n; i++ {
						xs[int8(rand.Int63())] = int8(rand.Int63())
					}
					ys := make(map[int8]int8)
					bin, err := ToBinary(&xs)
					Expect(err).ToNot(HaveOccurred())
					Expect(FromBinary(&ys, bin)).ToNot(HaveOccurred())
					return reflect.DeepEqual(xs, ys)
				}
				err := quick.Check(f, nil)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
