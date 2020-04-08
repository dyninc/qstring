package qstring

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestComparativeTimeParse(t *testing.T) {
	tme := "2006-01-02T15:04:05Z"
	testio := []struct {
		inp       string
		operator  string
		errString string
	}{
		{inp: tme, operator: "=", errString: ""},
		{inp: ">" + tme, operator: ">", errString: ""},
		{inp: "<" + tme, operator: "<", errString: ""},
		{inp: ">=" + tme, operator: ">=", errString: ""},
		{inp: "<=" + tme, operator: "<=", errString: ""},
		{inp: "<=" + tme, operator: "<=", errString: ""},
		{inp: "", operator: "=", errString: "qstring: Invalid Timestamp Query"},
		{inp: ">=", operator: "=", errString: "qstring: Invalid Timestamp Query"},
		{inp: ">=" + "foobar", operator: ">=",
			errString: `parsing time "foobar" as "2006-01-02T15:04:05Z07:00": cannot parse "foobar" as "2006"`},
	}

	var ct *ComparativeTime
	var err error
	for _, test := range testio {
		ct = NewComparativeTime()
		err = ct.Parse(test.inp)

		if ct.Operator != test.operator {
			t.Errorf("Expected operator %q, got %q", test.operator, ct.Operator)
		}

		if err == nil && len(test.errString) != 0 {
			t.Errorf("Expected error %q, got nil", test.errString)
		}

		if err != nil && err.Error() != test.errString {
			t.Errorf("Expected error %q, got %q", test.errString, err.Error())
		}
	}
}

func TestComparativeTimeUnmarshal(t *testing.T) {
	type Query struct {
		Created  ComparativeTime
		Modified ComparativeTime
	}

	createdTS := ">2006-01-02T15:04:05Z"
	updatedTS := "<=2016-01-02T15:04:05-07:00"

	query := url.Values{
		"created":  []string{createdTS},
		"modified": []string{updatedTS},
	}

	params := &Query{}
	err := Unmarshal(query, params)
	if err != nil {
		t.Fatal(err.Error())
	}

	created := params.Created.String()
	if created != createdTS {
		t.Errorf("Expected created ts of %s, got %s instead.", createdTS, created)
	}

	modified := params.Modified.String()
	if modified != updatedTS {
		t.Errorf("Expected update ts of %s, got %s instead.", updatedTS, modified)
	}
}

func TestComparativeTimeMarshalString(t *testing.T) {
	type Query struct {
		Created  ComparativeTime
		Modified ComparativeTime
	}

	createdTS := ">2006-01-02T15:04:05Z"
	created := NewComparativeTime()
	created.Parse(createdTS)
	updatedTS := "<=2016-01-02T15:04:05-07:00"
	updated := NewComparativeTime()
	updated.Parse(updatedTS)

	q := &Query{*created, *updated}
	result, err := MarshalString(q)
	if err != nil {
		t.Fatalf("Unable to marshal comparative timestamp: %s", err.Error())
	}

	var unescaped string
	unescaped, err = url.QueryUnescape(result)
	if err != nil {
		t.Fatalf("Unable to unescape query string %q: %q", result, err.Error())
	}
	expected := []string{"created=>2006-01-02T15:04:05Z",
		"modified=<=2016-01-02T15:04:05-07:00"}
	for _, ts := range expected {
		if !strings.Contains(unescaped, ts) {
			t.Errorf("Expected query string %s to contain %s", unescaped, ts)
		}
	}
}

func TestComparativeStringUnmarshal(t *testing.T) {
	type Query struct {
		Equals    ComparativeString
		Similar   ComparativeString
		Different ComparativeString
	}

	val1 := "stringValue1"
	equalsVal := fmt.Sprintf("%s", val1)
	val2 := "stringValue2"
	similarVal := fmt.Sprintf("~%s", val2)
	val3 := "stringValue3"
	diffVal := fmt.Sprintf("!%s", val3)

	query := url.Values{
		"equals":    []string{equalsVal},
		"different": []string{diffVal},
		"similar":   []string{similarVal},
	}

	params := &Query{}
	err := Unmarshal(query, params)
	if err != nil {
		t.Fatal(err.Error())
	}

	equals := params.Equals.String()
	if equals != equalsVal {
		t.Errorf("Expected equals val of %s, got %s instead.", equalsVal, equals)
	}
	similar := params.Similar.String()
	if similar != similarVal {
		t.Errorf("Expected similar val of %s, got %s instead.", similarVal, similar)
	}
	diff := params.Different.String()
	if diff != diffVal {
		t.Errorf("Expected different val of %s, got %s instead.", diffVal, diff)
	}
}

func TestComparativeStringMarshalString(t *testing.T) {
	type Query struct {
		Equals    ComparativeString
		Similar   ComparativeString
		Different ComparativeString
	}

	val1 := "stringValue1"
	equalsVal := fmt.Sprintf("%s", val1)
	equals := &ComparativeString{}
	equals.Parse(equalsVal)

	val2 := "stringValue2"
	similarVal := fmt.Sprintf("~%s", val2)
	similar := &ComparativeString{}
	similar.Parse(similarVal)

	val3 := "stringValue3"
	diffVal := fmt.Sprintf("!%s", val3)
	different := &ComparativeString{}
	different.Parse(diffVal)

	q := &Query{*equals, *similar, *different}
	result, err := MarshalString(q)
	if err != nil {
		t.Fatalf("Unable to marshal comparative timestamp: %s", err.Error())
	}

	var unescaped string
	unescaped, err = url.QueryUnescape(result)
	if err != nil {
		t.Fatalf("Unable to unescape query string %q: %q", result, err.Error())
	}
	expected := []string{
		fmt.Sprintf("different=!%s", val3),
		fmt.Sprintf("equals=%s", val1),
		fmt.Sprintf("similar=~%s", val2),
	}
	for _, ts := range expected {
		if !strings.Contains(unescaped, ts) {
			t.Errorf("Expected query string %s to contain %s", unescaped, ts)
		}
	}
}
