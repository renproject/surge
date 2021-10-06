package surge_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/surge"
	"github.com/renproject/surge/surgeutil"
)

var _ = Describe("String", func() {

	numTrials := 100

	ts := []reflect.Type{
		reflect.TypeOf(""),
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

func BenchmarkUnmarshalString(b *testing.B) {
	size := 20
	buf := make([]byte, size)
	buf, _, _ = surge.MarshalLen(uint32(size-surge.SizeHintU32), buf, size)
	var str string

	for i := 0; i < b.N; i++ {
		surge.UnmarshalString(&str, buf, size)
	}
}
