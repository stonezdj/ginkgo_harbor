package ginkgo_harbor_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGinkgoHarbor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinkgoHarbor Suite")
}
