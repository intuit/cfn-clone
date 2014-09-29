package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type describeStackResponse struct {
	Stacks []struct {
		Parameters []struct {
			ParameterKey   string
			ParameterValue string
		}
	}
}

func createStackCmd(name string, params map[string]string, template string) ([]string, error) {
	cmd := []string{
		"aws",
		"cloudformation",
		"create-stack",
		"--stack-name",
		name,
		"--template-body",
		"file:///" + template,
		"--capabilities \"CAPABILITY_IAM\"",
		"--parameters",
	}

	if err := noEchoParamsOverriden(params); err != nil {
		return []string{}, err
	}

	return append(cmd, cliParamsForCreate(params)...), nil
}

func createStack(name string, params map[string]string, template string) (string, error) {
	createCmd, err := createStackCmd(name, params, template)
	if err != nil {
		return "", err
	}

	fmt.Println("Going to run with command:")
	fmt.Printf("%s\n", strings.Join(createCmd, " "))

	cmd := exec.Command(createCmd[0], createCmd[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(output))
	}

	return string(output), nil
}

func newStackTemplateFile(sourceStack string, path string) (string, error) {
	t, err := template(sourceStack, path)
	if err != nil {
		return "", err
	}

	f, err := ioutil.TempFile("", "cfn-clone")
	if err != nil {
		fmt.Printf("Unable to create temp file for template. Error: %v", err)
		return "", err
	}

	defer f.Close()

	_, err = f.WriteString(t)
	if err != nil {
		fmt.Printf("Unable to write to temp file for template. Error: %v", err)
		return "", err
	}

	if err = f.Sync(); err != nil {
		fmt.Printf("Unable to flush write to temp file for template. Error: %v", err)
		return "", err
	}

	return f.Name(), nil
}

func noEchoParamsOverriden(params map[string]string) error {
	for k, v := range params {
		if v == "****" {
			return fmt.Errorf("NoEcho Paramater '%s' must have overrid value specified.", k)
		}
	}
	return nil
}

func cliParamsForCreate(params map[string]string) []string {
	p := []string{}
	for k, v := range params {
		escapedK := strings.Replace(k, ",", "\\,", -1)
		escapedV := strings.Replace(v, ",", "\\,", -1)
		p = append(p, "ParameterKey="+escapedK+",ParameterValue=\""+escapedV+"\"")
	}

	return p
}

func stackParametersCmd(stack string) []string {
	return []string{
		"aws",
		"cloudformation",
		"describe-stacks",
		"--stack-name",
		stack,
	}
}

func stackParameters(stack string) (map[string]string, error) {
	paramsCmd := stackParametersCmd(stack)

	cmd := exec.Command(paramsCmd[0], paramsCmd[1:]...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting attributes from source stack. Error: %s", string(output))
		os.Exit(1)
	}

	j := describeStackResponse{}
	if err = json.Unmarshal([]byte(string(output)), &j); err != nil {
		return map[string]string{}, err
	}

	params := map[string]string{}
	for _, p := range j.Stacks[0].Parameters {
		params[p.ParameterKey] = p.ParameterValue
	}

	return params, nil
}

func stackTemplateCmd(name string) []string {
	return []string{
		"aws",
		"cloudformation",
		"get-template",
		"--stack-name",
		name,
	}
}

func stackTemplate(name string) (string, error) {
	templateCmd := stackTemplateCmd(name)

	cmd := exec.Command(templateCmd[0], templateCmd[1:]...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	j := map[string]interface{}{}
	if err = json.Unmarshal([]byte(string(output)), &j); err != nil {
		return "", err
	}

	template, err := json.Marshal(j["TemplateBody"])
	if err != nil {
		return "", err
	}

	return string(template), nil
}

func template(sourceStack string, path string) (string, error) {
	if path == "" {
		return stackTemplate(sourceStack)
	} else {
		t, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Unable to read template file. Error: %v", err)
		}
		return string(t), nil
	}
}
