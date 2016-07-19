package qstring

import (
	"testing"
)

func TestTagParsing(t *testing.T) {
	testio := []struct {
		inp    string
		output string
		omit   bool
	}{
		{inp: "name,omitempty", output: "name", omit: true},
		{inp: "name", output: "name", omit: false},
		{inp: "name,", output: "name", omit: false},
		{inp: "name", output: "name", omit: false},
		{inp: "", output: "", omit: false},
		{inp: ",omitempty", output: "", omit: true},
		{inp: "-", output: "-", omit: false},
	}

	var name string
	var omit bool
	for _, test := range testio {
		name, omit = parseTag(test.inp)
		if name != test.output {
			t.Errorf("Expected tag name to be %q, got %q instead", test.output, name)
		}

		if omit != test.omit {
			t.Errorf("Expected omitempty to be %t, got %t instead", test.omit, omit)
		}
	}
}

func TestOmitEmpty(t *testing.T) {
	vis := 5
	testio := []struct {
		Visible     int
		Conditional int `qstring:"conditional,omitempty"`
		omit        bool
	}{
		{Visible: vis, Conditional: 5, omit: false},
		{Visible: vis, Conditional: 0, omit: true},
	}

	for _, test := range testio {
		values, _ := MarshalValues(&test)

		_, found := values["conditional"]
		if found && test.omit {
			t.Errorf("%d was unexpectedly marshaled", test.Conditional)
		} else if !found && !test.omit {
			t.Errorf("%d was not marshaled, but should have been", test.Conditional)
		}
	}
}
