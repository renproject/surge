package surge_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/surge/surgeutil"
)

var _ = Describe("Array", func() {

	numTrials := 10

	ts := []reflect.Type{
		reflect.TypeOf([0]int8{}),
		reflect.TypeOf([0]int16{}),
		reflect.TypeOf([0]int32{}),
		reflect.TypeOf([0]int64{}),
		reflect.TypeOf([0]uint8{}),
		reflect.TypeOf([0]uint16{}),
		reflect.TypeOf([0]uint32{}),
		reflect.TypeOf([0]uint64{}),
		reflect.TypeOf([0]bool{}),
		reflect.TypeOf([0]float32{}),
		reflect.TypeOf([0]float64{}),
		reflect.TypeOf([0]byte{}),
		reflect.TypeOf([0]string{}),
		reflect.TypeOf([1]int8{}),
		reflect.TypeOf([1]int16{}),
		reflect.TypeOf([1]int32{}),
		reflect.TypeOf([1]int64{}),
		reflect.TypeOf([1]uint8{}),
		reflect.TypeOf([1]uint16{}),
		reflect.TypeOf([1]uint32{}),
		reflect.TypeOf([1]uint64{}),
		reflect.TypeOf([1]bool{}),
		reflect.TypeOf([1]float32{}),
		reflect.TypeOf([1]float64{}),
		reflect.TypeOf([1]byte{}),
		reflect.TypeOf([1]string{}),
		reflect.TypeOf([100]int8{}),
		reflect.TypeOf([100]int16{}),
		reflect.TypeOf([100]int32{}),
		reflect.TypeOf([100]int64{}),
		reflect.TypeOf([100]uint8{}),
		reflect.TypeOf([100]uint16{}),
		reflect.TypeOf([100]uint32{}),
		reflect.TypeOf([100]uint64{}),
		reflect.TypeOf([100]bool{}),
		reflect.TypeOf([100]float32{}),
		reflect.TypeOf([100]float64{}),
		reflect.TypeOf([100]byte{}),
		reflect.TypeOf([100]string{}),
	}

	for _, t := range ts {
		t := t

		Context(fmt.Sprintf("when marshaling and then unmarshaling %v arrays", t), func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.MarshalUnmarshalCheck(t)).To(Succeed())
				}
			})
		})

		Context(fmt.Sprintf("when fuzzing %v arrays", t), func() {
			It("should not panic", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(func() { surgeutil.Fuzz(t) }).ToNot(Panic())
				}
			})
		})

		Context(fmt.Sprintf("when marshaling %v arrays", t), func() {
			Context("when the buffer is too small", func() {
				It("should return itself", func() {
					for trial := 0; trial < numTrials; trial++ {
						Expect(surgeutil.MarshalBufTooSmall(t)).To(Succeed())
					}
				})
			})

			Context("when the remaining memory quota is too small", func() {
				It("should return itself", func() {
					for trial := 0; trial < numTrials; trial++ {
						Expect(surgeutil.MarshalRemTooSmall(t)).To(Succeed())
					}
				})
			})
		})

		Context(fmt.Sprintf("when unmarshaling %v arrays", t), func() {
			Context("when the buffer is too small", func() {
				It("should return itself", func() {
					for trial := 0; trial < numTrials; trial++ {
						Expect(surgeutil.UnmarshalBufTooSmall(t)).To(Succeed())
					}
				})
			})

			Context("when the remaining memory quota is too small", func() {
				It("should return itself", func() {
					for trial := 0; trial < numTrials; trial++ {
						Expect(surgeutil.UnmarshalRemTooSmall(t)).To(Succeed())
					}
				})
			})
		})
	}
})
