package val_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestVal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Val Suite")
}
