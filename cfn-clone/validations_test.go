package main

import (
	"io/ioutil"
	"testing"
)

var cliParamsTcs = []struct {
	params        []string
	resultIsError bool
}{
	{[]string{"FOO=BAR", "BAZ=BLAH"}, false},
	{[]string{"FOO=BAR", "BAZ"}, true},
}

func TestValidateCliParameters(t *testing.T) {
	for _, tc := range cliParamsTcs {
		err := validateCliParameters(tc.params)
		if (err != nil) != tc.resultIsError {
			t.Fatalf("Expected '%v' got '%v' for '%v'", tc.resultIsError, err, tc.params)
		}
	}
}

var cliExistsTcs = []struct {
	cmd           string
	resultIsError bool
}{
	{"ls", false},
	{"no-way-this-exists", true},
}

func TestValidateCliExists(t *testing.T) {
	for _, tc := range cliExistsTcs {
		err := validateCliExists(tc.cmd)
		if (err != nil) != tc.resultIsError {
			t.Fatalf("Expected '%v' got '%v' for '%v'", tc.resultIsError, err, tc.cmd)
		}
	}
}

var templateExistsTcs = []struct {
	file          string
	createTmpFile bool
	resultIsError bool
}{
	{"", true, false},
	{"/no-way-this-exists", false, true},
}

func TestValidateTemplateExists(t *testing.T) {
	for _, tc := range templateExistsTcs {
		if tc.createTmpFile {
			f, err := ioutil.TempFile("", "cfn-clone-test")
			if err != nil {
				t.Fatalf("Unable to create temp file for testing ValidateTemplateExists")
			}

			defer f.Close()
			tc.file = f.Name()
		}

		err := validateTemplateExists(tc.file)
		if (err != nil) != tc.resultIsError {
			t.Fatalf("Expected '%v' got '%v' for '%v'", tc.resultIsError, err, tc.file)
		}
	}
}
