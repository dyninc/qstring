package qstring

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Marshaller defines the interface for performing custom marshalling of struct
// values into query strings
type Marshaller interface {
	MarshalQuery() (url.Values, error)
}

// Marshal marshals the provided struct into a url.Values collection
func Marshal(v interface{}) (url.Values, error) {
	var e encoder
	e.init(v)
	return e.marshal()
}

// Marshal marshals the provided struct into a raw query string and returns a
// conditional error
func MarshalString(v interface{}) (string, error) {
	vals, err := Marshal(v)
	if err != nil {
		return "", err
	}
	return vals.Encode(), nil
}

// An InvalidMarshalError describes an invalid argument passed to Marshal or
// MarshalValue. (The argument to Marshal must be a non-nil pointer.)
type InvalidMarshalError struct {
	Type reflect.Type
}

func (e InvalidMarshalError) Error() string {
	if e.Type == nil {
		return "qstring: MarshalString(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "qstring: MarshalString(non-pointer " + e.Type.String() + ")"
	}
	return "qstring: MarshalString(nil " + e.Type.String() + ")"
}

type encoder struct {
	data interface{}
}

func (e *encoder) init(v interface{}) *encoder {
	e.data = v
	return e
}

func (e *encoder) marshal() (url.Values, error) {
	rv := reflect.ValueOf(e.data)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return nil, &InvalidMarshalError{reflect.TypeOf(e.data)}
	}

	switch val := e.data.(type) {
	case Marshaller:
		return val.MarshalQuery()
	default:
		return e.value(rv)
	}
}

var textMarshallerElem = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
var stringerElem = reflect.TypeOf(new(fmt.Stringer)).Elem()

func (e *encoder) value(val reflect.Value) (url.Values, error) {
	elem := val.Elem()
	typ := elem.Type()

	var err error
	var output = make(url.Values)
	for i := 0; i < elem.NumField(); i++ {
		// pull out the qstring struct tag
		elemField := elem.Field(i)
		typField := typ.Field(i)
		qstring, omit := parseTag(typField.Tag.Get(Tag))
		if qstring == "" {
			// resolvable fields must have at least the `flag` struct tag
			qstring = strings.ToLower(typField.Name)
		}

		// determine if this is an unsettable field or was explicitly set to be
		// ignored
		if !elemField.CanSet() || qstring == "-" || (omit && isEmptyValue(elemField)) {
			continue
		}

		// verify if the element type implements compatible interfaces
		if val, ok := compatibleInterfaceValue(elemField); ok {
			output.Set(qstring, val)
			continue
		}
		// only do work if the current fields query string parameter was provided
		switch k := typField.Type.Kind(); k {
		case reflect.Slice:
			output[qstring] = marshalSlice(elemField)
		case reflect.Ptr:
			marshalStruct(output, qstring, reflect.Indirect(elemField), k)
		case reflect.Struct:
			marshalStruct(output, qstring, elemField, k)
		default:
			output.Set(qstring, marshalValue(elemField, k))
		}
	}
	return output, err
}

func marshalSlice(field reflect.Value) []string {
	var out []string
	for i := 0; i < field.Len(); i++ {
		out = append(out, marshalValue(field.Index(i), field.Index(i).Kind()))
	}
	return out
}

func compatibleInterfaceValue(field reflect.Value) (string, bool) {
	if field.Type().Implements(textMarshallerElem) {
		byt, _ := field.Interface().(encoding.TextMarshaler).MarshalText()
		return string(byt), true
	}
	if field.Type().Implements(stringerElem) {
		return field.Interface().(fmt.Stringer).String(), true
	}
	return "", false
}

func marshalValue(field reflect.Value, source reflect.Kind) string {
	if val, ok := compatibleInterfaceValue(field); ok {
		return val
	}
	switch source {
	case reflect.String:
		return field.String()
	case reflect.Bool:
		return strconv.FormatBool(field.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'G', -1, 64)
	case reflect.Struct:
		switch field.Interface().(type) {
		case encoding.TextMarshaler:
			byt, _ := field.Interface().(encoding.TextMarshaler).MarshalText()
			return string(byt)
		case time.Time:
			return field.Interface().(time.Time).Format(time.RFC3339)
		case ComparativeTime:
			return field.Interface().(ComparativeTime).String()
		}
	}
	return ""
}

func marshalStruct(output url.Values, qstring string, field reflect.Value, source reflect.Kind) error {
	var err error
	switch field.Interface().(type) {
	case time.Time, ComparativeTime:
		output.Set(qstring, marshalValue(field, source))
	default:
		var vals url.Values
		if field.CanAddr() {
			vals, err = Marshal(field.Addr().Interface())
		}

		if err != nil {
			return err
		}
		for key, list := range vals {
			if qstring != "" {
				key = qstring + "." + key
			}
			output[key] = list
		}
	}
	return nil
}
