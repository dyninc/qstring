package qstring

import (
	"reflect"
	"testing"
	"time"
)

func TestIsEmptyValue(t *testing.T) {
	var ts *TestStruct
	ts = nil
	testIO := []struct {
		inp      reflect.Value
		expected bool
	}{
		{inp: reflect.ValueOf([]int{0}), expected: false},
		{inp: reflect.ValueOf([]int{}), expected: true},
		{inp: reflect.ValueOf(map[string]int{"a": 0}), expected: false},
		{inp: reflect.ValueOf(map[string]int{}), expected: true},
		{inp: reflect.ValueOf(false), expected: true},
		{inp: reflect.ValueOf(true), expected: false},
		{inp: reflect.ValueOf(5), expected: false},
		{inp: reflect.ValueOf(0), expected: true},
		{inp: reflect.ValueOf(uint(5)), expected: false},
		{inp: reflect.ValueOf(uint(0)), expected: true},
		{inp: reflect.ValueOf(float32(5)), expected: false},
		{inp: reflect.ValueOf(float32(0)), expected: true},
		{inp: reflect.ValueOf(&TestStruct{}), expected: false},
		{inp: reflect.ValueOf(ts), expected: true},
		{inp: reflect.ValueOf(nil), expected: false},
		{inp: reflect.ValueOf(time.Time{}), expected: true},
		{inp: reflect.ValueOf(time.Now()), expected: false},
		{inp: reflect.ValueOf(*NewComparativeTime()), expected: true},
		{inp: reflect.ValueOf(ComparativeTime{Operator: "=", Time: time.Now()}),
			expected: false},
	}

	var result bool
	for _, test := range testIO {
		result = isEmptyValue(test.inp)
		if result != test.expected {
			t.Errorf("Expected %t for input %s, got %t", test.expected, test.inp, result)
		}
	}
}
