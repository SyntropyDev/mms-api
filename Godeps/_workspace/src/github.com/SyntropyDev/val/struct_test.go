package val_test

import (
	"github.com/SyntropyDev/val"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testStruct1 struct {
	Required string `val:"nonzero"`
}

type testStruct2 struct {
	Email    string `val:"email"`
	HexColor string `val:"hexcolor"`
}

type testStruct3 struct {
	Url      string `val:"url"`
	Username string `val:"alphanum | minlen(5) | maxlen(10)"`
}

type testStruct4 struct {
	Int   int     `val:"gt(18)"`
	Int2  int64   `val:"gte(100)"`
	Float float64 `val:"lt(19.9)"`
	Uint  uint    `val:"lte(18)"`
}

type testStruct5 struct {
	IntSlice []int `val:"each(gt(18))"`
}

type testStruct6 struct {
	Lat float64 `val:"lat"`
	Lon float64 `val:"lon"`
}

var validStructs = []interface{}{
	testStruct1{"a"},
	testStruct2{"loganjspears@gmail.com", "#ffffff"},
	testStruct2{"loganjspears@gmail.com", "#FFFFFF"},
	testStruct3{"http://www.google.com", "logan12345"},
	testStruct4{19, 400, -10.0000, 18},
	testStruct5{[]int{19, 20, 21}},
	testStruct6{0.0, 0.0},
	testStruct6{90.0, -180.0},
}

var invalidStructs = []interface{}{
	testStruct1{Required: ""},
	testStruct2{"loganjspears@gmail", "#ffffff"},
	testStruct2{"loganjspears@gmail.com", "#fhffff"},
	testStruct3{"http://google", "logan12345"},
	testStruct3{"http://www.google.com", "log1"},
	testStruct3{"http://www.google.com", "logan100101001100101001"},
	testStruct4{14, 400, -10.0000, 18},
	testStruct4{19, 99, -10.0000, 18},
	testStruct4{19, 400, 19.90000001, 18},
	testStruct4{19, 400, -10.0000, 19},
	testStruct5{[]int{19, 20, 10}},
	testStruct6{90.1, -180.0},
	testStruct6{90.0, -180.1},
}

var _ = Describe("Struct", func() {

	Describe("Valid", func() {

		It("should be valid for the values", func() {
			for _, v := range validStructs {
				_, errs := val.Struct(v)
				Expect(errs).To(BeEmpty())
			}
		})
	})

	Describe("Invalid", func() {

		It("should be invalid for the values", func() {
			for _, v := range invalidStructs {
				valid, _ := val.Struct(v)
				Expect(valid).To(BeFalse())
			}
		})
	})
})
