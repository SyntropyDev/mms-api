package merge_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMerge(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Merge Suite")
}
