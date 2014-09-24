Val
===

[![Build Status](https://drone.io/github.com/SyntropyDev/val/status.png)](https://drone.io/github.com/SyntropyDev/val/latest)

A struct validation library for go.  Val took code and inspiration from:

- https://github.com/astaxie/beego/tree/master/validation
- https://github.com/onsi/gomega

```go
type ValStruct struct {
	Required        string            `val:”nonzero”`
	Email           string            `val:”email”`
	Website         string            `val:"url"`
	IpAddress       string            `val:"ip"`
	PetName         string            `val:"alpha"`
	Phone           string            `val:"num | len(10)"`
	Password        string            `val:”alphanum | minlen(5) | maxlen(15)"`
	DayOfWeek       int               `val:”gte(0) | lt(7)”`
	Lat             float64           `val:”lat”`
	Lng             float64           `val:”lon”`
	TennisGameScore string            `val:”in(love,15,30,40)”`
	NewSuperHero    string            `val:”notin(Superman,Batman,The Flash)”`
	UserInfo        map[string]string `val:"haskey(id)"`
	Zipcode         string            `val:”match(^\d{5}(?:[-\s]\d{4})?$)”`
	Color           string            `val:”hexcolor”`
	Ages            []int             `val:”each( gt(18) | lt(35) )”`
	Explode         string            `val:”panic(nonzero)”`
}
```
