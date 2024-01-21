package repository

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repositories Suite")
}
