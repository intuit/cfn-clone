package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCliParamsForCreate(t *testing.T) {
	expected := []string{
		"ParameterKey=foo,ParameterValue=\"bar\\,1\"",
		"ParameterKey=bar,ParameterValue=\"2\\,baz\"",
	}
	params := map[string]string{"foo": "bar,1", "bar": "2,baz"}
	result := cliParamsForCreate(params)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("Expected '%v' got '%v'", expected, result)
	}
}

func TestCreateStackCmd(t *testing.T) {
	name := "foo"
	params := map[string]string{"param1": "val1,valx", "param2": "val2,valy"}
	template := "/var/tmp/new_template.json"

	expected := []string{
		"aws",
		"cloudformation",
		"create-stack",
		"--output",
		"json",
		"--stack-name",
		name,
		"--template-body",
		"file:///" + template,
		"--capabilities",
		"CAPABILITY_IAM",
		"--parameters",
		"ParameterKey=param1,ParameterValue=\"val1\\,valx\"",
		"ParameterKey=param2,ParameterValue=\"val2\\,valy\"",
	}

	cmd, _ := createStackCmd(name, params, template)

	// because order of maps are not guaranteed
	if !reflect.DeepEqual(cmd[:8], expected[:8]) {
		t.Fatalf("Expected '%s' got '%s'", expected[:8], cmd[:8])
	}

	for i := 8; i < 10; i++ {
		if expected[i] != cmd[8] && expected[i] != cmd[9] {
			t.Fatalf("Missing parameter. Expected to find '%s' in '%s' or '%s'", expected[i], cmd[8], cmd[9])
		}
	}
}

func TestStackParamsCmd(t *testing.T) {
	name := "foo"

	expected := []string{
		"aws",
		"cloudformation",
		"describe-stacks",
		"--output",
		"json",
		"--stack-name",
		name,
	}

	cmd := stackParametersCmd(name)

	if !reflect.DeepEqual(cmd, expected) {
		t.Fatalf("Expected '%s' got '%s'", expected, cmd)
	}
}

func TestStackTemplateCmd(t *testing.T) {
	name := "foo"

	expected := []string{
		"aws",
		"cloudformation",
		"get-template",
		"--output",
		"json",
		"--stack-name",
		name,
	}

	cmd := stackTemplateCmd(name)

	if !reflect.DeepEqual(cmd, expected) {
		t.Fatalf("Expected '%s' got '%s'", expected, cmd)
	}
}

func TestTemplate(t *testing.T) {
	s := `{"foo": "bar"}`
	f, err := ioutil.TempFile("", "cfn-clone-test")
	if err != nil {
		t.Fatalf("Unable to create temp file for testing ValidateTemplateExists")
	}

	defer os.Remove(f.Name())
	defer f.Close()

	f.WriteString(s)
	f.Sync()
	f.Close()

	data, _ := template("", f.Name())

	if data != s {
		t.Fatalf("Expected '%s' got '%s'", s, data)
	}
}
