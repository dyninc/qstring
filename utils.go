package qstring

import (
	"reflect"
	"strings"
	"time"
)

// isEmptyValue returns true if the provided reflect.Value
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		switch t := v.Interface().(type) {
		case time.Time:
			return t.IsZero()
		case ComparativeTime:
			return t.Time.IsZero()
		}
	}
	return false
}

// parseTag splits a struct field's qstring tag into its name and, if an
// optional omitempty option was provided, a boolean indicating this is
// returned and a boolean
func parseTag(tag string) (name string, omitempty bool, comma bool) {
	elements := strings.Split(tag, ",")
	if len(elements) == 1 {
		return tag, false, false
	}

	name = elements[0]
	for _, opt := range elements[1:] {
		if opt == "omitempty" {
			omitempty = true
		}
		if opt == "comma" {
			comma = true
		}
	}
	return
}
