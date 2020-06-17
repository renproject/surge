package surge_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("String", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x string) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHintString(x)+excess)
						rem := surge.SizeHintString(x) + excess

						tail, tailRem, err := surge.MarshalString(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					f := func(x string) bool {
						size := surge.SizeHintString(x)
						buf := make([]byte, size)
						for rem := 0; rem < size; rem++ {
							_, _, err := surge.MarshalString(x, buf, rem)
							Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
						}
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})
		})

		Context("when the buffer is not big enough", func() {
			It("should return an error", func() {
				f := func(x string) bool {
					size := surge.SizeHintString(x)
					rem := size
					for b := 0; b < size; b++ {
						buf := make([]byte, b)
						_, _, err := surge.MarshalString(x, buf, rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					}
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

					x := ""
					Expect(func() { surge.UnmarshalString(&x, buf, rem) }).ToNot(Panic())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should return the original value", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x string) bool {
						excess := r.Int() % 100
						rem := surge.SizeHintString(x) + excess
						buf := make([]byte, surge.SizeHintString(x))
						_, _, err := surge.MarshalString(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						y := ""
						_, tailRem, err := surge.UnmarshalString(&y, buf[:], rem)
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
					f := func(x string) bool {
						rem := surge.SizeHintString(x)
						buf := make([]byte, surge.SizeHintString(x))
						_, _, err := surge.MarshalString(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())
						rem = surge.SizeHintString(x) - 1

						y := ""
						_, _, err = surge.UnmarshalString(&y, buf[:], rem)
						Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))

						y = ""
						_, _, err = surge.UnmarshalString(&y, buf[:], 1)
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
				f := func(x string) bool {
					excess := r.Int() % 100
					rem := surge.SizeHintString(x) + excess
					buf := make([]byte, surge.SizeHintString(x))
					_, _, err := surge.MarshalString(x, buf, rem)
					Expect(err).ToNot(HaveOccurred())
					buf = buf[1:]

					y := ""
					_, _, err = surge.UnmarshalString(&y, buf[:], rem)
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
