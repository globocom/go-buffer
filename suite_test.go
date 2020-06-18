package buffer_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGoBuffer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-buffer suite")
}
