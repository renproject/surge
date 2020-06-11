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
						buf := make([]byte, surge.SizeHintBool+excess)
						rem := surge.SizeHintBool + excess

						tail, tailRem, err := surge.MarshalBool(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
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
						buf := make([]byte, surge.SizeHintBool+excess)
						rem := surge.SizeHintBool - 1

						tail, tailRem, err := surge.MarshalBool(x, buf, rem)
						Expect(tail).To(HaveLen(len(buf)))
						Expect(tailRem).To(Equal(rem))
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
					buf := make([]byte, surge.SizeHintBool-1)
					rem := surge.SizeHintBool + excess

					tail, tailRem, err := surge.MarshalBool(x, buf, rem)
					Expect(tail).To(HaveLen(len(buf)))
					Expect(tailRem).To(Equal(rem))
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
					buf := make([]byte, surge.SizeHintBool+excess)
					rem := surge.SizeHintBool + excess

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
						rem := surge.SizeHintBool + excess

						x := false
						tail, tailRem, err := surge.UnmarshalBool(&x, buf[:], rem)
						Expect(tail).To(HaveLen(len(buf) - surge.SizeHintBool))
						Expect(tailRem).To(Equal(excess))
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
						rem := surge.SizeHintBool - 1
						x := false
						tail, tailRem, err := surge.UnmarshalBool(&x, buf[:], rem)
						Expect(tail).To(HaveLen(len(buf)))
						Expect(tailRem).To(Equal(rem))
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
				f := func(buf [surge.SizeHintBool - 1]byte) bool {
					excess := r.Int() % 100
					rem := surge.SizeHintBool + excess

					x := false
					tail, tailRem, err := surge.UnmarshalBool(&x, buf[:], rem)
					Expect(tail).To(HaveLen(len(buf)))
					Expect(tailRem).To(Equal(rem))
					Expect(err).To(Equal(surge.ErrUnexpectedEndOfBuffer))
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})
	})
})
