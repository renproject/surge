package surge_test

import (
	"fmt"
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/surge/surgeutil"
)

var _ = Describe("Map", func() {

	numTrials := 10

	ts := []reflect.Type{
		reflect.TypeOf(map[int8]int8{}),
		reflect.TypeOf(map[int16]int16{}),
		reflect.TypeOf(map[int32]int32{}),
		reflect.TypeOf(map[int64]int64{}),
		reflect.TypeOf(map[uint8]uint8{}),
		reflect.TypeOf(map[uint16]uint16{}),
		reflect.TypeOf(map[uint32]uint32{}),
		reflect.TypeOf(map[uint64]uint64{}),
		reflect.TypeOf(map[bool]bool{}),
		reflect.TypeOf(map[float32]float32{}),
		reflect.TypeOf(map[float64]float64{}),
		reflect.TypeOf(map[byte]byte{}),
		reflect.TypeOf(map[string]string{}),
		reflect.TypeOf(map[string]int8{}),
		reflect.TypeOf(map[string]int16{}),
		reflect.TypeOf(map[string]int32{}),
		reflect.TypeOf(map[string]int64{}),
		reflect.TypeOf(map[string]uint8{}),
		reflect.TypeOf(map[string]uint16{}),
		reflect.TypeOf(map[string]uint32{}),
		reflect.TypeOf(map[string]uint64{}),
		reflect.TypeOf(map[string]bool{}),
		reflect.TypeOf(map[string]float32{}),
		reflect.TypeOf(map[string]float64{}),
		reflect.TypeOf(map[string]byte{}),
		reflect.TypeOf(map[string]string{}),
	}

	for _, t := range ts {
		Context(fmt.Sprintf("when marshaling and then unmarshaling %v maps", t), func() {
			It("should return itself", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(surgeutil.MarshalUnmarshalCheck(t)).To(Succeed())
				}
			})
		})

		Context(fmt.Sprintf("when fuzzing %v maps", t), func() {
			It("should not panic", func() {
				for trial := 0; trial < numTrials; trial++ {
					Expect(func() { surgeutil.Fuzz(t) }).ToNot(Panic())
				}
			})
		})

		Context(fmt.Sprintf("when marshaling %v maps", t), func() {
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

		Context(fmt.Sprintf("when unmarshaling %v maps", t), func() {
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
