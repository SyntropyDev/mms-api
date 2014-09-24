package val

import (
	"errors"
	"github.com/onsi/gomega/matchers"
	"reflect"
)

const (
	regexEmail    = `^([a-z0-9_\.-]+)@([\da-z\.-]+)\.([a-z\.]{2,6})$`
	regexHexColor = `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	regexUrl      = `^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	regexIp       = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	regexNum      = `^[1-9]\d*(\.\d+)?$`
	regexAlpha    = `^[a-zA-Z]*$`
)

type matcher interface {
	Match(actual interface{}) (success bool, err error)
}

type ValFunc func(k string, v interface{}) error

var Nonzero = func(k string, v interface{}) error {
	m := &matchers.BeZeroMatcher{}
	zero, _ := m.Match(v)
	if zero {
		return formatError(k)
	}
	return nil
}

func Matches(regex string) ValFunc {
	m := &matchers.MatchRegexpMatcher{Regexp: regex}
	return match(m)
}

var Email = Matches(regexEmail)
var HexColor = Matches(regexHexColor)
var Url = Matches(regexUrl)
var Ip = Matches(regexIp)
var Alpha = Matches(regexAlpha)
var Num = Matches(regexNum)
var AlphaNum = combineValFuncs(Matches("[a-zA-Z]+"), Matches("[0-9]+"))

func HasKey(k interface{}) ValFunc {
	m := &matchers.HaveKeyMatcher{Key: k}
	return match(m)
}

func Gt(v interface{}) ValFunc {
	return numericalMatch(">", v)
}

func Gte(v interface{}) ValFunc {
	return numericalMatch(">=", v)
}

func Lt(v interface{}) ValFunc {
	return numericalMatch("<", v)
}

func Lte(v interface{}) ValFunc {
	return numericalMatch("<=", v)
}

var Lat = combineValFuncs(
	numericalMatch("<=", 90.0),
	numericalMatch(">=", -90.0))

var Lon = combineValFuncs(
	numericalMatch("<=", 180.0),
	numericalMatch(">=", -180.0))

func In(list []interface{}) ValFunc {
	return func(k string, v interface{}) error {
		in := false
		for _, e := range list {
			in = in || reflect.DeepEqual(e, v)
		}
		if !in {
			return formatError(k)
		}
		return nil
	}
}

func NotIn(list []interface{}) ValFunc {
	return func(k string, v interface{}) error {
		for _, e := range list {
			if reflect.DeepEqual(e, v) {
				return formatError(k)
			}
		}
		return nil
	}
}

func Len(l int) ValFunc {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length != l {
			return formatError(k)
		}
		return nil
	}
}

func MinLen(l int) ValFunc {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length < l {
			return formatError(k)
		}
		return nil
	}
}

func MaxLen(l int) ValFunc {
	return func(k string, v interface{}) error {
		length, ok := lengthOf(v)
		if !ok || length > l {
			return formatError(k)
		}
		return nil
	}
}

func Each(valFuncs ...ValFunc) ValFunc {
	return func(k string, v interface{}) error {
		if !isArrayOrSlice(v) {
			return formatError(k)
		}
		value := reflect.ValueOf(v)
		for i := 0; i < value.Len(); i++ {
			iFace := value.Index(i).Interface()
			for _, f := range valFuncs {
				if err := f(k, iFace); err != nil {
					return formatError(k)
				}
			}
		}
		return nil
	}
}

func Panic(valFuncs ...ValFunc) ValFunc {
	return func(k string, v interface{}) error {
		for _, f := range valFuncs {
			if err := f(k, v); err != nil {
				panic(err)
			}
		}
		return nil
	}
}

func numericalMatch(comparator string, v interface{}) ValFunc {
	m := &matchers.BeNumericallyMatcher{
		Comparator: comparator,
		CompareTo:  []interface{}{v},
	}
	return match(m)
}

func match(m matcher) ValFunc {
	return func(k string, v interface{}) error {
		matches, _ := m.Match(v)
		if !matches {
			return formatError(k)
		}
		return nil
	}
}

func combineValFuncs(valFuncs ...ValFunc) ValFunc {
	return func(k string, v interface{}) error {
		for _, f := range valFuncs {
			if err := f(k, v); err != nil {
				return err
			}
		}
		return nil
	}
}

func formatError(k string) error {
	return errors.New(k + " did not pass validation.")
}
