package surge_test

import (
	"github.com/renproject/surge"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	Context("when creating an unsupported type error", func() {
		It("should contain the name of the unsupported type", func() {
			err := surge.NewErrUnsupportedMarshalType(float64(0))
			Expect(err.Error()).To(ContainSubstring("unsupported"))
			Expect(err.Error()).To(ContainSubstring("marshal"))
			Expect(err.Error()).To(ContainSubstring("float64"))

			err = surge.NewErrUnsupportedUnmarshalType(float64(0))
			Expect(err.Error()).To(ContainSubstring("unsupported"))
			Expect(err.Error()).To(ContainSubstring("unmarshal"))
			Expect(err.Error()).To(ContainSubstring("float64"))
		})
	})
})
