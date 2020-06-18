package surge_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Float", func() {

	numTrials := 100

	ts := []reflect.Type{
		reflect.TypeOf(float32(0)),
		reflect.TypeOf(float64(0)),
	}

	Context("when marshaling and then unmarshaling", func() {
		It("should return itself", func() {
			for trial := 0; trial < numTrials; trial++ {
				for _, t := range ts {
					Expect(MarshalUnmarshalCheck(t)).To(Succeed())
				}
			}
		})
	})

	Context("when fuzzing", func() {
		It("should not panic", func() {
			for trial := 0; trial < numTrials; trial++ {
				for _, t := range ts {
					Expect(func() { FuzzCheck(t) }).ToNot(Panic())
				}
			}
		})
	})

	Context("when marshaling", func() {
		Context("when the buffer is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(MarshalBufTooSmall(t)).To(Succeed())
					}
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(MarshalRemTooSmall(t)).To(Succeed())
					}
				}
			})
		})
	})

	Context("when unmarshaling", func() {
		Context("when the buffer is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(UnmarshalBufTooSmall(t)).To(Succeed())
					}
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(UnmarshalRemTooSmall(t)).To(Succeed())
					}
				}
			})
		})
	})
})
