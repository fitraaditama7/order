package order

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Order Service Suite")
}
