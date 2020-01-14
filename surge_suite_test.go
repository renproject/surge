package surge_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSurge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Surge Suite")
}
