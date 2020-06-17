package surge_test

import (
	"math/rand"
	"testing/quick"
	"time"

	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Float32", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float32) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHintF32+excess)
						rem := surge.SizeHintF32 + excess

						tail, tailRem, err := surge.MarshalF32(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						y := float32(0)
						_, _, err = surge.UnmarshalF32(&y, buf, rem)
						Expect(err).ToNot(HaveOccurred())
						Expect(x).To(Equal(y))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float32) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHintF32+excess)
						rem := surge.SizeHintF32 - 1

						tail, tailRem, err := surge.MarshalF32(x, buf, rem)
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
				f := func(x float32) bool {
					excess := r.Int() % 100
					buf := make([]byte, surge.SizeHintF32-1)
					rem := surge.SizeHintF32 + excess

					tail, tailRem, err := surge.MarshalF32(x, buf, rem)
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
					buf := make([]byte, surge.SizeHintF32+excess)
					rem := surge.SizeHintF32 + excess

					x := float32(0)
					Expect(func() { surge.UnmarshalF32(&x, buf, rem) }).ToNot(Panic())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should return the original value", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float32) bool {
						excess := r.Int() % 100
						rem := surge.SizeHintF32 + excess
						buf := make([]byte, surge.SizeHintF32)
						_, _, err := surge.MarshalF32(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						y := float32(0)
						tail, tailRem, err := surge.UnmarshalF32(&y, buf[:], rem)
						Expect(tail).To(HaveLen(len(buf) - surge.SizeHintF32))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						Expect(x).To(Equal(y))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					f := func(buf [100]byte) bool {
						rem := surge.SizeHintF32 - 1
						x := float32(0)
						tail, tailRem, err := surge.UnmarshalF32(&x, buf[:], rem)
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
				f := func(buf [surge.SizeHintF32 - 1]byte) bool {
					excess := r.Int() % 100
					rem := surge.SizeHintF32 + excess

					x := float32(0)
					tail, tailRem, err := surge.UnmarshalF32(&x, buf[:], rem)
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

var _ = Describe("Float64", func() {
	Context("when marshaling", func() {
		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should not return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float64) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHintF64+excess)
						rem := surge.SizeHintF64 + excess

						tail, tailRem, err := surge.MarshalF64(x, buf, rem)
						Expect(tail).To(HaveLen(excess))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						y := float64(0)
						_, _, err = surge.UnmarshalF64(&y, buf, rem)
						Expect(err).ToNot(HaveOccurred())
						Expect(x).To(Equal(y))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float64) bool {
						excess := r.Int() % 100
						buf := make([]byte, surge.SizeHintF64+excess)
						rem := surge.SizeHintF64 - 1

						tail, tailRem, err := surge.MarshalF64(x, buf, rem)
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
				f := func(x float64) bool {
					excess := r.Int() % 100
					buf := make([]byte, surge.SizeHintF64-1)
					rem := surge.SizeHintF64 + excess

					tail, tailRem, err := surge.MarshalF64(x, buf, rem)
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
					buf := make([]byte, surge.SizeHintF64+excess)
					rem := surge.SizeHintF64 + excess

					x := float64(0)
					Expect(func() { surge.UnmarshalF64(&x, buf, rem) }).ToNot(Panic())
					return true
				}
				Expect(quick.Check(f, nil)).To(Succeed())
			})
		})

		Context("when the buffer is big enough", func() {
			Context("when there are sufficient remaining bytes", func() {
				It("should return the original value", func() {
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					f := func(x float64) bool {
						excess := r.Int() % 100
						rem := surge.SizeHintF64 + excess
						buf := make([]byte, surge.SizeHintF64)
						_, _, err := surge.MarshalF64(x, buf, rem)
						Expect(err).ToNot(HaveOccurred())

						y := float64(0)
						tail, tailRem, err := surge.UnmarshalF64(&y, buf[:], rem)
						Expect(tail).To(HaveLen(len(buf) - surge.SizeHintF64))
						Expect(tailRem).To(Equal(excess))
						Expect(err).ToNot(HaveOccurred())

						Expect(x).To(Equal(y))
						return true
					}
					Expect(quick.Check(f, nil)).To(Succeed())
				})
			})

			Context("when there are not sufficient remaining bytes", func() {
				It("should return an error", func() {
					f := func(buf [100]byte) bool {
						rem := surge.SizeHintF64 - 1
						x := float64(0)
						tail, tailRem, err := surge.UnmarshalF64(&x, buf[:], rem)
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
				f := func(buf [surge.SizeHintF64 - 1]byte) bool {
					excess := r.Int() % 100
					rem := surge.SizeHintF64 + excess

					x := float64(0)
					tail, tailRem, err := surge.UnmarshalF64(&x, buf[:], rem)
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
