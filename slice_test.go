package surge_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Slice", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x []float32) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHint(x)+excess)
						rem := surge.SizeHint(x) + excess

						tail, tailRem, err := surge.Marshal(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				Context("when there are less remaining bytes than the length of the array", func() {
					It("should not return an error", func() {
						r := rand.New(rand.NewSource(time.Now().UnixNano()))
						f := func(x []float64) bool {
							excess := r.Int() % 100
							buf := make([]byte, surge.SizeHint(x)+excess)
							rem := surge.SizeHint(x) - 1

							_, tailRem, err := surge.Marshal(x, buf, rem)
							Expect(tailRem).To(BeNumerically("<=", rem))
							Expect(err).To(HaveOccurred())

							return true
						}
						Expect(quick.Check(f, nil)).To(Succeed())
					})
				})

				It("should return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x []int8) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHint(x)+excess)
						rem := surge.SizeHint(x) - 1

						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})
		})

		Context("when the buffer is not big enough", func() {
			Context("when the buffer is shorter than the array", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x []int16) bool {
						excess := r.Int() % 100
						buf := make([]byte, len(x))
						rem := surge.SizeHint(x) + excess

						_, tailRem, err := surge.Marshal(x, buf, rem)
						Expect(tailRem).To(BeNumerically("<=", rem))
						Expect(err).To(HaveOccurred())

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			It("should return an error", func() {
				r := rand.New(rand.NewSource(time.Now().UnixNano()))
				f := func(x []int32) bool {
					excess := r.Int() % 100
					buf := make([]byte, surge.SizeHint(x)-1)
					rem := surge.SizeHint(x) + excess

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

					x := []int64{}
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
					f := func(x []uint8) bool {
						excess := r.Int() % 100
						rem := 2*surge.SizeHint(x) + excess
						buf := make([]byte, surge.SizeHint(x))
						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						y := []uint8{}
						_, tailRem, err := surge.Unmarshal(&y, buf[:], rem)
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
					f := func(x []uint16) bool {
						rem := surge.SizeHint(x)
						buf := make([]byte, surge.SizeHint(x))
						_, _, err := surge.Marshal(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())
						rem = surge.SizeHint(x) - 1

						y := []uint16{}
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
				f := func(x []uint32) bool {
					excess := r.Int() % 100
					rem := surge.SizeHint(x) + excess
					buf := make([]byte, surge.SizeHint(x))
					_, _, err := surge.Marshal(x, buf, rem)
					Expect(err).ToNot(HaveOccurred())
					buf = buf[1:]

					y := []uint32{}
					_, _, err = surge.Unmarshal(&y, buf[:], rem)
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
