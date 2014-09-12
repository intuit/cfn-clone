package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	flags "github.com/jessevdk/go-flags"
)

type DescribeStackResponse struct {
	Stacks []struct {
		Parameters []struct {
			ParameterKey string
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
		p = append(p, "ParameterKey=" + k + ",ParameterValue=" + v + " ")
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

func verifyAwsCliExists() {
	_, err := exec.LookPath("aws")
	if err != nil {
		fmt.Printf("Unable to find the AWS CLI in your PATH")
		os.Exit(1)
	}
}

func verifyCliParameters(params []string) {
	invalid := false
	for _, p := range params {
		v := strings.SplitN(p, "=", 2)
		if len(v) != 2 {
			fmt.Printf("Attribute '%s' must be '=' separated key, value", p)
			invalid = true
		}
	}

	if invalid {
		fmt.Println("")
		os.Exit(1)
	}
}

func verifyTemplateExists(path string) {
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("Unable to find the template '%s'. Is that the correct path?", path)
			os.Exit(1)
		}
	}
}

func verifySourceStackExists(name string) {
	args := []string{
		"cloudformation",
		"describe-stacks",
		"--stack-name",
		name,
	}

	cmd := exec.Command("aws", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error verifying source stack. Error: %s", string(output))
		os.Exit(1)
	}
}

func main() {
	var opts struct {
		Attributes []string `short:"a" long:"attributes" description:"'=' separated attribute and value"`
		NewName string `short:"n" long:"new-name" description:"Name for new stack" required:"true"`
		SourceName string `short:"s" long:"source-name" description:"Name of source stack to clone" required:"true"`
		Template string `short:"t" long:"template" description:"Path to a new template file"`
	}

	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()

	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	verifyAwsCliExists()
	verifyTemplateExists(opts.Template)
	verifySourceStackExists(opts.SourceName)
	verifyCliParameters(opts.Attributes)

	newTemplate, err := newStackTemplateFile(opts.SourceName, opts.Template)
	if err != nil {
		fmt.Printf("Erroring getting the template for cloning. Error: %v", err)
		os.Exit(1)
	}
	defer os.Remove(newTemplate)

	parameters, err := sourceStackParameters(opts.SourceName)
	if err != nil {
		fmt.Printf("Error getting source stack parameters. Error: %v", err)
		os.Exit(1)
	}

	for k, v := range cliParameters(opts.Attributes) {
		parameters[k] = v
	}

	fmt.Printf("Our merged parameters are %v\n", parameters)

	fmt.Println("Going to clone")

	createNewStack(opts.NewName, parameters, newTemplate)

	os.Exit(0)
}
