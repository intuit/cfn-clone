package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type DescribeStackResponse struct {
	Stacks []struct {
		Parameters []struct {
			ParameterKey   string
			ParameterValue string
		}
	}
}

func cliParameters(attribs []string) map[string]string {
	parameters := map[string]string{}
	for _, a := range attribs {
		p := strings.SplitN(a, "=", 2)
		parameters[p[0]] = p[1]
	}

	return parameters
}

func createNewStack(name string, params map[string]string, template string) {
	args := []string{
		"cloudformation",
		"create-stack",
		"--stack-name",
		name,
		"--template-body",
		"file:///" + template,
		"--parameters",
	}

	args = append(args, paramsForCreate(params)...)

	fmt.Println("Going to run with command:")
	fmt.Printf("aws %s\n", strings.Join(args, " "))

	cmd := exec.Command("aws", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error creating new stack. Error: %s", string(output))
		os.Exit(1)
	}

	fmt.Printf("%s", string(output))
}

func paramsForCreate(params map[string]string) []string {
	p := []string{}
	for k, v := range params {
		p = append(p, "ParameterKey="+k+",ParameterValue="+v+" ")
	}

	return p
}

func sourceStackParameters(stack string) (map[string]string, error) {
	args := []string{
		"cloudformation",
		"describe-stacks",
		"--stack-name",
		stack,
	}

	cmd := exec.Command("aws", args...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error getting attributes from source stack. Error: %s", string(output))
		os.Exit(1)
	}

	j := DescribeStackResponse{}
	if err = json.Unmarshal([]byte(string(output)), &j); err != nil {
		return map[string]string{}, err
	}

	params := map[string]string{}
	for _, p := range j.Stacks[0].Parameters {
		params[p.ParameterKey] = p.ParameterValue
	}

	return params, nil
}

func sourceStackTemplate(name string) (string, error) {
	args := []string{
		"cloudformation",
		"get-template",
		"--stack-name",
		name,
	}

	cmd := exec.Command("aws", args...)

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
		return sourceStackTemplate(sourceStack)
	} else {
		t, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Printf("Unable to read template file. Error: %v", err)
		}
		return string(t), nil
	}
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

func main() {
	options := parseCliArgs()

	newTemplate, err := newStackTemplateFile(options.SourceName, options.Template)
	if err != nil {
		fmt.Printf("Erroring getting the template for cloning. Error: %v", err)
		os.Exit(1)
	}
	defer os.Remove(newTemplate)

	parameters, err := sourceStackParameters(options.SourceName)
	if err != nil {
		fmt.Printf("Error getting source stack parameters. Error: %v", err)
		os.Exit(1)
	}

	for k, v := range cliParameters(options.Attributes) {
		parameters[k] = v
	}

	fmt.Printf("Our merged parameters are %v\n", parameters)

	fmt.Println("Going to clone")

	createNewStack(options.NewName, parameters, newTemplate)

	os.Exit(0)
}
