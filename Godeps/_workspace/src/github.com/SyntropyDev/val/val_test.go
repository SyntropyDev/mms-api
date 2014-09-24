package val_test

import (
	"github.com/SyntropyDev/val"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestStruct1 struct {
	Value string
}

var _ = Describe("Val", func() {

	Describe("Validator", func() {

		Describe("Nonzero", func() {

			Describe("Valid", func() {

				It("should be valid for 1", func() {
					v := val.New()
					v.Add("test", 1, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid for abc", func() {
					v := val.New()
					v.Add("test", "abc", val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid for struct", func() {
					v := val.New()
					v.Add("test", TestStruct1{"a"}, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid for ptr to struct", func() {
					v := val.New()
					v.Add("test", &TestStruct1{"a"}, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})
			})

			Describe("Invalid", func() {

				It("should be invalid for 0", func() {
					v := val.New()
					v.Add("test", 0, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for \"\"", func() {
					v := val.New()
					v.Add("test", "", val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be valid for empty struct", func() {
					v := val.New()
					v.Add("test", TestStruct1{}, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be valid for nil", func() {
					v := val.New()
					v.Add("test", nil, val.Nonzero)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})
			})

		})

		Describe("Email", func() {

			Describe("Valid", func() {

				It("should be valid logan@syntropy.io", func() {
					v := val.New()
					v.Add("test", "logan@syntropy.io", val.Email)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid for loganjspears@gmail.com", func() {
					v := val.New()
					v.Add("test", "loganjspears@gmail.com", val.Email)
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

			})

			Describe("Invalid", func() {

				It("should be invalid logansyntropy.io", func() {
					v := val.New()
					v.Add("test", "logansyntropy.io", val.Email)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for loganjspears@gmail", func() {
					v := val.New()
					v.Add("test", "loganjspears@gmail", val.Email)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for @yahoo.com", func() {
					v := val.New()
					v.Add("test", "@yahoo.com", val.Email)
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})
			})

		})

		Describe("val.Matches", func() {

			Describe("Valid", func() {

				It("should be valid for regex", func() {
					v := val.New()
					v.Add("test", "lmy-us3r_n4m3", val.Matches("^[a-z0-9_-]{3,16}$"))
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid for regex", func() {
					v := val.New()
					v.Add("test", "myp4ssw0rd", val.Matches("^[a-z0-9_-]{6,18}$"))
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

			})

			Describe("Invalid", func() {

				It("should be invalid for regex", func() {
					v := val.New()
					v.Add("test", "th1s1s-wayt00_l0ngt0beausername", val.Matches("^[a-z0-9_-]{3,16}$"))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for regex", func() {
					v := val.New()
					v.Add("test", "mypa$$w0rd", val.Matches("^[a-z0-9_-]{6,18}$"))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for integer", func() {
					v := val.New()
					v.Add("test", 5, val.Matches("^[a-z0-9_-]{6,18}$"))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})
			})

		})

		Describe("Panic", func() {

			Describe("Valid", func() {

				It("should not panic if condition is met", func() {
					v := val.New()
					v.Add("test", "1", val.Panic(val.MaxLen(2)))
					f := func() {
						v.Validate()
					}
					Expect(f).ShouldNot(Panic())
				})

			})

			Describe("Invalid", func() {

				It("should panic if condition is not met", func() {
					v := val.New()
					v.Add("test", "13333", val.Panic(val.MaxLen(2)))
					f := func() {
						v.Validate()
					}
					Expect(f).Should(Panic())
				})
			})

		})

		Describe("Gt", func() {

			Describe("Valid", func() {

				It("should be valid int greater than the value", func() {
					v := val.New()
					v.Add("test", 2, val.Gt(1))
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

				It("should be valid float greater than the value", func() {
					v := val.New()
					v.Add("test", 5.0, val.Gt(4.9999999999999))
					valid, _ := v.Validate()
					Expect(valid).To(BeTrue())
				})

			})

			Describe("Invalid", func() {

				It("should be invalid for int less than the value", func() {
					v := val.New()
					v.Add("test", 1, val.Gt(2))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for float less than the value", func() {
					v := val.New()
					v.Add("test", 4.9999999999999, val.Gt(5.0))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})

				It("should be invalid for string", func() {
					v := val.New()
					v.Add("test", "2", val.Gt(1))
					valid, _ := v.Validate()
					Expect(valid).To(BeFalse())
				})
			})
		})
	})

	Describe("Lt", func() {

		Describe("Valid", func() {

			It("should be valid int less than the value", func() {
				v := val.New()
				v.Add("test", 1, val.Lt(2))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid float less than the value", func() {
				v := val.New()
				v.Add("test", 4.9999999999999, val.Lt(5.0))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

		})

		Describe("Invalid", func() {

			It("should be invalid for int greater than the value", func() {
				v := val.New()
				v.Add("test", 2, val.Lt(1))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for float greater than the value", func() {
				v := val.New()
				v.Add("test", 5.0, val.Lt(4.9999999999999))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for string", func() {
				v := val.New()
				v.Add("test", "2", val.Lt(3))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})
	})

	Describe("Lat", func() {

		Describe("Valid", func() {

			It("should be valid for valid lat", func() {
				v := val.New()
				v.Add("test", 1, val.Lat)
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for valid lat", func() {
				v := val.New()
				v.Add("test", 90.0, val.Lat)
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

		})

		Describe("Invalid", func() {

			It("should be invalid for invalid lat", func() {
				v := val.New()
				v.Add("test", -90.1, val.Lat)
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for invalid lat", func() {
				v := val.New()
				v.Add("test", 91.0, val.Lat)
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

		})
	})

	Describe("Lon", func() {

		Describe("Valid", func() {

			It("should be valid for valid lon", func() {
				v := val.New()
				v.Add("test", 1, val.Lon)
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for valid lon", func() {
				v := val.New()
				v.Add("test", 180.0, val.Lon)
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

		})

		Describe("Invalid", func() {

			It("should be invalid for invalid lon", func() {
				v := val.New()
				v.Add("test", -181.0, val.Lon)
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for invalid lon", func() {
				v := val.New()
				v.Add("test", 180.01, val.Lon)
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

		})
	})

	Describe("In", func() {

		Describe("Valid", func() {

			It("should be valid for element in slice", func() {
				v := val.New()
				v.Add("test", "in", val.In([]interface{}{"in", "out"}))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for element in slice", func() {
				v := val.New()
				v.Add("test", 2, val.In([]interface{}{1, 2}))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

		})

		Describe("Invalid", func() {

			It("should be invalid for element not in slice", func() {
				v := val.New()
				v.Add("test", "not in", val.In([]interface{}{"in", "out"}))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be valid for element not in slice", func() {
				v := val.New()
				v.Add("test", 3, val.In([]interface{}{1, 2}))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for list of the wrong type", func() {
				v := val.New()
				v.Add("test", "a", val.In([]interface{}{1, 2}))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})

	})

	Describe("NotIn", func() {

		Describe("Valid", func() {

			It("should be valid for element not in slice", func() {
				v := val.New()
				v.Add("test", "not in", val.NotIn([]interface{}{"in", "out"}))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for element not in slice", func() {
				v := val.New()
				v.Add("test", 3, val.NotIn([]interface{}{1, 2}))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

		})

		Describe("Invalid", func() {

			It("should be invalid for element in slice", func() {
				v := val.New()
				v.Add("test", "in", val.NotIn([]interface{}{"in", "out"}))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for element in slice", func() {
				v := val.New()
				v.Add("test", 2, val.NotIn([]interface{}{1, 2}))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})

	})

	Describe("Len", func() {

		Describe("Valid", func() {

			It("should be valid for a string of the correct length", func() {
				v := val.New()
				v.Add("test", "123456", val.Len(6))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a slice of the correct length", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.Len(3))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})
		})

		Describe("Invalid", func() {

			It("should be invalid for a string longer than value", func() {
				v := val.New()
				v.Add("test", "123456", val.Len(7))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for a slice shorter than value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.Len(4))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})

	})

	Describe("MinLen", func() {

		Describe("Valid", func() {

			It("should be valid for a string longer than value", func() {
				v := val.New()
				v.Add("test", "123456", val.MinLen(5))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a string as long as the value", func() {
				v := val.New()
				v.Add("test", "123456", val.MinLen(6))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a slice longer than value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MinLen(2))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a slice as long as the value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MinLen(3))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})
		})

		Describe("Invalid", func() {

			It("should be invalid for a string shorter than value", func() {
				v := val.New()
				v.Add("test", "123456", val.MinLen(7))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for a slice shorter than value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MinLen(4))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})
	})

	Describe("MaxLen", func() {

		Describe("Valid", func() {

			It("should be valid for a string shorter than value", func() {
				v := val.New()
				v.Add("test", "123456", val.MaxLen(7))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a string as long as the value", func() {
				v := val.New()
				v.Add("test", "123456", val.MaxLen(6))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a slice shorter than value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MaxLen(4))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for a slice as long as the value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MaxLen(3))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})
		})

		Describe("Invalid", func() {

			It("should be invalid for a string longer than value", func() {
				v := val.New()
				v.Add("test", "123456", val.MaxLen(4))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be invalid for a slice longer than value", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.MaxLen(2))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})
	})

	Describe("Each", func() {

		Describe("Valid", func() {

			It("should be valid for each value in the list", func() {
				v := val.New()
				v.Add("test", []string{"1", "12", "123"}, val.Each(val.MaxLen(3)))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})

			It("should be valid for each value in the list", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.Each(val.Gt(0)))
				valid, _ := v.Validate()
				Expect(valid).To(BeTrue())
			})
		})

		Describe("Invalid", func() {

			It("should be invalid if a value in the list isn't valid", func() {
				v := val.New()
				v.Add("test", []string{"1", "12", "123"}, val.Each(val.MaxLen(2)))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})

			It("should be valid for each value in the list", func() {
				v := val.New()
				v.Add("test", []int{1, 2, 3}, val.Each(val.Gt(1)))
				valid, _ := v.Validate()
				Expect(valid).To(BeFalse())
			})
		})
	})
})
