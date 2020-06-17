package surge_test

import (
	"bytes"
	"math/rand"
	"testing/quick"
	"time"
	"unsafe"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x map[int8]int16) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHint(x)+excess)
						rem := surge.SizeHint(x) + 48*len(x) + excess

						tail, tailRem, err := surge.Marshal(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})

				It("should return the same bytes for multiple marshals", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x map[string]map[string]string) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHint(x)+excess)
						rem := 2*surge.SizeHint(x) + excess + 48*len(x)

						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						for i := 0; i < 10; i++ {
							buf2 := make([]byte, surge.SizeHint(x)+excess)
							rem2 := 2*surge.SizeHint(x) + excess + 48*len(x)

							_, _, err := surge.Marshal(x, buf2, rem2)
							Expect(err).ToNot(HaveOccurred())
							Expect(bytes.Equal(buf, buf2)).To(BeTrue())
						}

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x map[string]string) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHint(x)+excess)

						rem := surge.SizeHint(x) + 48*len(x) - 1
						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						rem = surge.SizeHint(x) - 1
						_, _, err = surge.Marshal(x, buf, rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						rem = 48*len(x) - 1
						_, _, err = surge.Marshal(x, buf, rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})
		})

		Context("when the buffer is not big enough", func() {
			It("should return an error", func() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				f := func(x map[string]int32) bool {
					excess := r.Int() % 100
					buf := make([]byte, surge.SizeHint(x)-1)
					rem := 2*surge.SizeHint(x) + excess + 48*len(x)

					_, _, err := surge.Marshal(x, buf, rem)
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})

	Context("when unmarshaling", func() {
		Context("when fuzzing", func() {
			It("should not panic", func() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				f := func(data []byte) bool {
					excess := r.Int() % 100
					buf := make([]byte, len(data)+excess)
					rem := len(data) + excess

					x := map[string][100]uint64{}
					Expect(func() { surge.Unmarshal(&x, buf, rem) }).ToNot(Panic())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should return the original value", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x map[int8][100]int64) bool {
						excess := r.Int() % 100
						rem := surge.SizeHint(x) + 48*len(x) + excess
						buf := make([]byte, rem)
						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						y := map[int8][100]int64{}
						tail, tailRem, err := surge.Unmarshal(&y, buf[:], rem)
						Expect(tail).To(HaveLen(len(buf) - surge.SizeHint(x)))
						Expect(tailRem).To(BeNumerically("<=", rem))
						Expect(err).ToNot(HaveOccurred())

						Expect(x).To(Equal(y))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					f := func(x map[int16]int32) bool {
						rem := surge.SizeHint(x) + 48*len(x)
						buf := make([]byte, surge.SizeHint(x))
						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						rem = surge.SizeHint(x) + int(unsafe.Sizeof(x)) - 1
						y := map[int16]int32{}
						_, _, err = surge.Unmarshal(&y, buf[:], rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						rem = surge.SizeHint(x) - 1
						y = map[int16]int32{}
						_, _, err = surge.Unmarshal(&y, buf[:], rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						rem = int(unsafe.Sizeof(x)) - 1
						y = map[int16]int32{}
						_, _, err = surge.Unmarshal(&y, buf[:], rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})
		})

		Context("when the buffer is not big enough", func() {
			It("should return an error", func() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				f := func(x map[string]string) bool {
					excess := r.Int() % 100
					rem := 2*surge.SizeHint(x) + excess + 48*len(x)
					buf := make([]byte, surge.SizeHint(x))
					_, _, err := surge.Marshal(x, buf, rem)
					Expect(err).ToNot(HaveOccurred())
					buf = buf[1:]

					y := map[string]string{}
					_, _, err = surge.Unmarshal(&y, buf[:], rem)
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
