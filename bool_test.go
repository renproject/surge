package surge_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Boolean", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x bool) bool {
						excess := r.Int() % 100
						buf := make([]byte, 1+excess)
						rem := 1 + excess

						tail, rem, err := surge.MarshalBool(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(rem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						if x {
							Expect(buf[0]).To(Equal(uint8(1)))
						} else {
							Expect(buf[0]).To(Equal(uint8(0)))
						}
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x bool) bool {
						excess := r.Int() % 100
						buf := make([]byte, 1+excess)
						rem := 0

						tail, rem, err := surge.MarshalBool(x, buf, rem)
						Expect(tail).To(HaveLen(1 + excess))
						Expect(rem).To(Equal(0))
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
				f := func(x bool) bool {
					excess := r.Int() % 100
					buf := make([]byte, 0)
					rem := 1 + excess

					tail, rem, err := surge.MarshalBool(x, buf, rem)
					Expect(tail).To(HaveLen(0))
					Expect(rem).To(Equal(1 + excess))
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
					buf := make([]byte, 1+excess)
					rem := 1 + excess

					x := false
					Expect(func() { surge.UnmarshalBool(&x, buf, rem) }).ToNot(Panic())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(buf [100]byte) bool {
						excess := r.Int() % 100
						rem := 1 + excess

						x := false
						tail, rem, err := surge.UnmarshalBool(&x, buf[:], rem)
						Expect(tail).To(HaveLen(99))
						Expect(rem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						if buf[0] == 0 {
							Expect(x).To(BeFalse())
						} else {
							Expect(x).To(BeTrue())
						}
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					f := func(buf [100]byte) bool {
						rem := 0
						x := false
						tail, rem, err := surge.UnmarshalBool(&x, buf[:], rem)
						Expect(tail).To(HaveLen(100))
						Expect(rem).To(Equal(0))
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
				f := func() bool {
					excess := r.Int() % 100
					rem := 1 + excess

					x := false
					tail, rem, err := surge.UnmarshalBool(&x, []byte{}, rem)
					Expect(tail).To(HaveLen(0))
					Expect(rem).To(Equal(1 + excess))
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
