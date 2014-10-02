package merge_test

import (
	"github.com/SyntropyDev/merge"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestStruct struct {
	One   string `merge:"true"`
	Two   int    `merge:"false"`
	Three string
}

var _ = Describe("Merge", func() {

	Describe("Validaiton", func() {

		It("should panic if to is not a pointer to a struct", func() {
			from := &TestStruct{"from", 1, "a"}
			to := TestStruct{"to", 0, "b"}
			f := func() {
				merge.Wl(from, to, "Two", "Three")
			}
			Expect(f).To(Panic())
		})
	})

	Describe("Wl", func() {

		It("should merge the struct whitelisted fields", func() {
			from := &TestStruct{"from", 1, "a"}
			to := &TestStruct{"to", 0, "b"}
			merge.Wl(from, to, "Two", "Three")
			Expect(to.One).To(Equal("to"))
			Expect(to.Two).To(Equal(1))
			Expect(to.Three).To(Equal("a"))
		})
	})

	Describe("TagWl", func() {

		It("should merge the struct whitelisted fields", func() {
			from := &TestStruct{"from", 1, "a"}
			to := &TestStruct{"to", 0, "b"}
			merge.TagWl(from, to)
			Expect(to.One).To(Equal("from"))
			Expect(to.Two).To(Equal(0))
			Expect(to.Three).To(Equal("b"))
		})
	})

	Describe("Bl", func() {

		It("should not merge the struct blacklisted fields", func() {
			from := &TestStruct{"from", 1, "a"}
			to := &TestStruct{"to", 0, "b"}
			merge.Bl(from, to, "Two")
			Expect(to.One).To(Equal("from"))
			Expect(to.Two).To(Equal(0))
			Expect(to.Three).To(Equal("a"))
		})
	})

	Describe("TagBl", func() {

		It("should not merge the struct blacklisted fields", func() {
			from := &TestStruct{"from", 1, "a"}
			to := &TestStruct{"to", 0, "b"}
			merge.TagBl(from, to)
			Expect(to.One).To(Equal("from"))
			Expect(to.Two).To(Equal(0))
			Expect(to.Three).To(Equal("a"))
		})
	})
})
