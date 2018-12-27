package stringutils

import (
	"testing"
)

func TestIn(t *testing.T) {
	v := []string{"Parrot", "Cat", "Camel", "Dog", "Gopher"}

	tests := []struct {
		i bool
		t string
		e bool
	}{
		{false, "dog", false},
		{false, "Dog", true},
		{false, "Pigeon", false},

		{true, "dog", true},
		{true, "Dog", true},
		{true, "Pigeon", false},
	}

	for i, test := range tests {
		r := In(test.i, test.t, v...)
		if r != test.e {
			t.Fatalf("[FAIL] Test %d: Expected %q with insensitive set to %t to return %t, but got %t instead", i+1, test.t, test.i, test.e, r)
		}
	}
}

func TestRandomString(t *testing.T) {
	tests := []int{10, 100, 1000}
	for i, test := range tests {
		s := RandomString(test)
		r := len(s)
		if r != test {
			t.Fatalf("[FAIL] Test %d: Expected length of %q to be %d, but got %d instead.", i+1, s, test, r)
		}
	}
}

func TestTitledString_String(t *testing.T) {
	var tests = []struct {
		s        string
		expected string
	}{
		{"test", "Test"},
		{"this-is-a-test", "This is a test"},
	}

	for i, test := range tests {
		r := TitledString(test.s)
		if r.String() != test.expected {
			t.Fatalf("%d: %s Expected %q, but got %q.", i, test.s, test.expected, r)
		}
		t.Logf("%d: %s Expected %q.", i, test.s, test.expected)
	}
}

func TestReplaceBetween(t *testing.T) {
	os := "This is a **test** string"
	r := os
	e := "This is a [test] string"
	ReplaceBetween(&r, "**", "[", "]")
	if r != e {
		t.Fatalf("%s: Expected %q, but got %q.", os, e, r)
	}
	t.Logf("%s: Expected %q.", os, e)
}
