package surge_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/surge/surgeutil"
)

var _ = Describe("Slice", func() {

	numTrials := 10

	ts := []reflect.Type{
		reflect.TypeOf([]int8{}),
		reflect.TypeOf([]int16{}),
		reflect.TypeOf([]int32{}),
		reflect.TypeOf([]int64{}),
		reflect.TypeOf([]uint8{}),
		reflect.TypeOf([]uint16{}),
		reflect.TypeOf([]uint32{}),
		reflect.TypeOf([]uint64{}),
		reflect.TypeOf([]bool{}),
		reflect.TypeOf([]float32{}),
		reflect.TypeOf([]float64{}),
		reflect.TypeOf([]byte{}),
		reflect.TypeOf([]string{}),
	}

	Context("when marshaling and then unmarshaling", func() {
		It("should return itself", func() {
			for trial := 0; trial < numTrials; trial++ {
				for _, t := range ts {
					Expect(surgeutil.MarshalUnmarshalCheck(t)).To(Succeed())
				}
			}
		})
	})

	Context("when fuzzing", func() {
		It("should not panic", func() {
			for trial := 0; trial < numTrials; trial++ {
				for _, t := range ts {
					Expect(func() { surgeutil.Fuzz(t) }).ToNot(Panic())
				}
			}
		})
	})

	Context("when marshaling", func() {
		Context("when the buffer is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(surgeutil.MarshalBufTooSmall(t)).To(Succeed())
					}
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(surgeutil.MarshalRemTooSmall(t)).To(Succeed())
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
						Expect(surgeutil.UnmarshalBufTooSmall(t)).To(Succeed())
					}
				}
			})
		})

		Context("when the remaining memory quota is too small", func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					for _, t := range ts {
						Expect(surgeutil.UnmarshalRemTooSmall(t)).To(Succeed())
					}
				}
			})
		})
	})
})
