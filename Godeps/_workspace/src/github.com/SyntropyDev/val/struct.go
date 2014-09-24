package val

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	structTagKey = "val"
)

type valFuncTag struct {
	tagKey   string
	convFunc func(string) ValFunc
}

var (
	nonzeroTag = valFuncTag{
		tagKey: "nonzero",
		convFunc: func(s string) ValFunc {
			return Nonzero
		},
	}
	emailTag = valFuncTag{
		tagKey: "email",
		convFunc: func(s string) ValFunc {
			return Email
		},
	}
	hexColorTag = valFuncTag{
		tagKey: "hexcolor",
		convFunc: func(s string) ValFunc {
			return HexColor
		},
	}
	urlTag = valFuncTag{
		tagKey: "url",
		convFunc: func(s string) ValFunc {
			return Url
		},
	}
	ipTag = valFuncTag{
		tagKey: "ip",
		convFunc: func(s string) ValFunc {
			return Ip
		},
	}
	alphaTag = valFuncTag{
		tagKey: "alpha",
		convFunc: func(s string) ValFunc {
			return Alpha
		},
	}
	numTag = valFuncTag{
		tagKey: "num",
		convFunc: func(s string) ValFunc {
			return Num
		},
	}
	alphaNumTag = valFuncTag{
		tagKey: "alphanum",
		convFunc: func(s string) ValFunc {
			return AlphaNum
		},
	}
	gtTag = valFuncTag{
		tagKey: "gt",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic("val gt: " + err.Error())
			}
			return Gt(n)
		},
	}
	gteTag = valFuncTag{
		tagKey: "gte",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic("val gte: " + err.Error())
			}
			return Gte(n)
		},
	}
	ltTag = valFuncTag{
		tagKey: "lt",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic("val lt: " + err.Error())
			}
			return Lt(n)
		},
	}
	lteTag = valFuncTag{
		tagKey: "lte",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic("val lte: " + err.Error())
			}
			return Lte(n)
		},
	}
	latTag = valFuncTag{
		tagKey: "lat",
		convFunc: func(s string) ValFunc {
			return Lat
		},
	}
	lonTag = valFuncTag{
		tagKey: "lon",
		convFunc: func(s string) ValFunc {
			return Lon
		},
	}
	inTag = valFuncTag{
		tagKey: "in",
		convFunc: func(s string) ValFunc {
			list := strings.Split(s, ",")
			iList := []interface{}{}
			for _, i := range list {
				iList = append(iList, i)
			}
			return In(iList)
		},
	}
	notInTag = valFuncTag{
		tagKey: "notin",
		convFunc: func(s string) ValFunc {
			list := strings.Split(s, ",")
			iList := []interface{}{}
			for _, i := range list {
				iList = append(iList, i)
			}
			return NotIn(iList)
		},
	}
	hasKeyTag = valFuncTag{
		tagKey: "haskey",
		convFunc: func(s string) ValFunc {
			return HasKey(s)
		},
	}
	matchesTag = valFuncTag{
		tagKey: "matches",
		convFunc: func(s string) ValFunc {
			return Matches(s)
		},
	}
	lenTag = valFuncTag{
		tagKey: "len",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic("val len: " + err.Error())
			}
			return Len(int(n))
		},
	}
	minLenTag = valFuncTag{
		tagKey: "minlen",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic("val minlen: " + err.Error())
			}
			return MinLen(int(n))
		},
	}
	maxLenTag = valFuncTag{
		tagKey: "maxlen",
		convFunc: func(s string) ValFunc {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic("val maxlen: " + err.Error())
			}
			return MaxLen(int(n))
		},
	}
	eachTag = valFuncTag{
		tagKey: "each",
		convFunc: func(s string) ValFunc {
			return Each(valFuncsFromTag(s)...)
		},
	}
	panicTag = valFuncTag{
		tagKey: "panic",
		convFunc: func(s string) ValFunc {
			return Panic(valFuncsFromTag(s)...)
		},
	}

	valFuncTags []valFuncTag
)

func init() {
	valFuncTags = []valFuncTag{nonzeroTag, emailTag, hexColorTag, urlTag,
		ipTag, alphaTag, numTag, alphaNumTag, gtTag, gteTag, ltTag, lteTag,
		latTag, lonTag, inTag, notInTag, matchesTag, lenTag, minLenTag,
		maxLenTag, eachTag, panicTag}
}

func Struct(str interface{}) (bool, map[string]error) {
	objT := reflect.TypeOf(str)
	objV := reflect.ValueOf(str)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		panic("str must be a struct or a pointer to a struct")
	}

	v := New()
	for i := 0; i < objT.NumField(); i++ {
		field := objT.Field(i)
		tag := field.Tag.Get(structTagKey)
		name := field.Name
		funcs := valFuncsFromTag(tag)

		value := objV.Field(i).Interface()
		for _, f := range funcs {
			v.Add(name, value, f)
		}
	}

	return v.Validate()
}

func valFuncsFromTag(tag string) []ValFunc {
	funcs := []ValFunc{}
	sects := strings.Split(tag, "|")
	tSects := []string{}
	for _, s := range sects {
		tSects = append(tSects, strings.TrimSpace(s))
	}
	for _, s := range tSects {
		nCap := nonCaptureString(s)
		cap := captureString(s)
		for _, t := range valFuncTags {
			if nCap == t.tagKey {
				f := t.convFunc(cap)
				funcs = append(funcs, f)
				break
			}
		}
	}
	return funcs
}

func nonCaptureString(s string) string {
	start := strings.Index(s, "(")
	if start == -1 {
		return s
	}
	return strings.TrimSpace(s[0:start])
}

func captureString(s string) string {
	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	if start == -1 || end == -1 || start > end {
		return ""
	}
	return strings.TrimSpace(s[start+1 : end])
}
