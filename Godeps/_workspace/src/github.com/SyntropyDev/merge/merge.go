package merge

import (
	"reflect"
	"strings"
)

type ListType int

const (
	structTagKey = "merge"
	tagTrue      = "true"
	tagFalse     = "false"
)

func Wl(from interface{}, to interface{}, keys ...string) {
	merge(from, to, keys...)
}

func Bl(from interface{}, to interface{}, keys ...string) {
	structKeys := keysFromStruct(from)
	wlKeys := []string{}
	for _, k := range structKeys {
		found := false
		for _, blk := range keys {
			if blk == k {
				found = true
			}
		}
		if !found {
			wlKeys = append(wlKeys, k)
		}
	}
	merge(from, to, wlKeys...)
}

func TagWl(from interface{}, to interface{}) {
	fType, _ := formReflect(from)
	keys := keysFromTags(fType, wlParser)
	Wl(from, to, keys...)
}

func TagBl(from interface{}, to interface{}) {
	fType, _ := formReflect(from)
	keys := keysFromTags(fType, blParser)
	Wl(from, to, keys...)
}

func merge(from interface{}, to interface{}, keys ...string) {
	if !isStructPtr(reflect.TypeOf(to)) {
		panic("to must be a pointer to a struct")
	}
	fromType, fromValue := formReflect(from)
	toType, toValue := formReflect(to)

	if fromType != toType {
		panic("from and to parameters must be the same type")
	}

	for _, k := range keys {
		fv := fromValue.FieldByName(k)
		tv := toValue.FieldByName(k)

		if tv.IsValid() && tv.CanSet() {
			tv.Set(fv)
		}
	}
}

func keysFromTags(t reflect.Type, f func(string) bool) []string {
	keys := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(structTagKey)
		if f(tag) {
			keys = append(keys, field.Name)
		}
	}
	return keys
}

func keysFromStruct(i interface{}) []string {
	fType, _ := formReflect(i)
	keys := []string{}

	for i := 0; i < fType.NumField(); i++ {
		field := fType.Field(i)
		keys = append(keys, field.Name)
	}

	return keys
}

func wlParser(tag string) bool {
	tag = strings.TrimSpace(tag)
	return tag == tagTrue
}

func blParser(tag string) bool {
	tag = strings.TrimSpace(tag)
	return tag != tagFalse
}

func formReflect(i interface{}) (reflect.Type, reflect.Value) {
	objT := reflect.TypeOf(i)
	objV := reflect.ValueOf(i)
	switch {
	case isStruct(objT):
	case isStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		panic("str must be a struct or a pointer to a struct")
	}
	return objT, objV
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
