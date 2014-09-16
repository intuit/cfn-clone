package main

import (
	"reflect"
	"testing"
)

var paramsFromCliTcs = []struct {
	attribs []string
	result  map[string]string
}{
	{[]string{"FOO=BAR"}, map[string]string{"FOO": "BAR"}},
	{[]string{"FOO=BAR=BAZ"}, map[string]string{"FOO": "BAR=BAZ"}},
}

func TestParamsFromCli(t *testing.T) {
	for _, tc := range paramsFromCliTcs {
		p := paramsFromCli(tc.attribs)

		if !reflect.DeepEqual(tc.result, p) {
			t.Fatalf("Expected '%v' got '%v'", tc.result, p)
		}
	}
}
