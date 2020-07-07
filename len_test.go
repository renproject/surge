package surge_test

import (
	"math/rand"
	"unsafe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/surge"
)

var _ = Describe("Slice length", func() {
	trials := 100

	Context("when marshalling and then unmarshalling", func() {
		var x, y uint32
		var bs [4]byte

		It("should return itself", func() {
			for i := 0; i < trials; i++ {
				// Generate
				x = rand.Uint32() & 0x000f_ffff                // x < 2^20
				size := rand.Intn(127) + 1                     // size < 2^7
				rem := surge.SizeHintU32 + int(x*uint32(size)) // rem < 2^27 + 4
				// Marshal
				_, _, err := surge.MarshalLen(x, bs[:], 4)
				Expect(err).ToNot(HaveOccurred())
				// Unmarshal
				_, _, err = surge.UnmarshalLen(&y, size, bs[:], rem)
				Expect(err).ToNot(HaveOccurred())
				// Equality
				Expect(y).To(Equal(x))
			}
		})
	})

	Context("when unmarshaling", func() {
		It("should return an error when rem is too small", func() {
			const maxLen uint32 = 1000

			var x, y uint32
			buf := make([]byte, maxLen)

			for i := 0; i < trials; i++ {
				// Generate value
				x = rand.Uint32() % maxLen
				size := rand.Intn(127) + 1
				// Marshal the value so that we can attempt to unmarshal the resulting data
				_, _, err := surge.MarshalLen(x, buf, surge.MaxBytes)
				Expect(err).ToNot(HaveOccurred())
				// Unmarshal with rem that is too small
				c := int(x)*size + surge.SizeHintU32
				for rem := 0; rem < c; rem++ {
					_, _, err := surge.UnmarshalLen(&y, size, buf, rem)
					Expect(err).To(HaveOccurred())
				}
			}
		})

		It("should return an error when there is an overflow", func() {
			// Overflow can only occur on systems where the int type has 64
			// bits.
			if unsafe.Sizeof(int(0)) != 64 {
				Succeed()
			}

			var y uint32
			var bs [4]byte

			x := ^uint32(0)
			size := 1 << 33

			_, _, err := surge.MarshalLen(x, bs[:], surge.MaxBytes)
			Expect(err).ToNot(HaveOccurred())

			_, _, err = surge.UnmarshalLen(&y, size, bs[:], surge.MaxBytes)
			Expect(err).To(Equal(surge.ErrLengthOverflow))
		})
	})
})
