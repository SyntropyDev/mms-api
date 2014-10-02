merge
=====

Go library for merging structs and maps

[![Build Status](https://drone.io/github.com/SyntropyDev/merge/status.png)](https://drone.io/github.com/SyntropyDev/merge/latest)

Usage:
```go

type StructToMerge struct {
	One   string
	Two   int    `merge:"false"`
	Three string `merge:"true"`
}

func main() {
	from := &StructToMerge{"from", 1, "b"}
	to := &StructToMerge{"to", 0, "c"}

	// whitelist fields
	merge.Wl(from, to, "Three")
	// to = {to,0,b}

	from = &StructToMerge{"from", 1, "b"}
	to = &StructToMerge{"to", 0, "c"}

	// blacklist fields
	merge.Bl(from, to, "Two")
	// to = {from,0,b}

	from = &StructToMerge{"from", 1, "b"}
	to = &StructToMerge{"to", 0, "c"}

	// whitelist using struct tags
	merge.TagWl(from, to)
	// to = {to,0,b}

	from = &StructToMerge{"from", 1, "b"}
	to = &StructToMerge{"to", 0, "c"}

	// blacklist using struct tags
	merge.TagBl(from, to)
	// to = {from,0,b}
}


```
