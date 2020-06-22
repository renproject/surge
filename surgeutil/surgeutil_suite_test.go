package surgeutil_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSurgeutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Surgeutil Suite")
}
