package main

import (
	"regexp"
	"testing"
)

func TestPrettyParameters(t *testing.T) {
	in := map[string]string{
		"foo":                 "bar",
		"apple":               "banana",
		"this-is-much-longer": "yep",
	}

	pattern := `The merged parameters are:\napple\s+banana\nfoo\s+bar\nthis-is-much-longer\s+yep\n`
	expected := regexp.MustCompile(pattern)

	out := prettyParameters(in)

	if !expected.MatchString(out) {
		t.Fatalf("Expected '%v' got '%v'", expected, out)
	}
}
