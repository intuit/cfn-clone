package main

import (
	"fmt"
	"os"
	"strings"
)

func cliParameters(attribs []string) map[string]string {
	parameters := map[string]string{}
	for _, a := range attribs {
		p := strings.SplitN(a, "=", 2)
		parameters[p[0]] = p[1]
	}

	return parameters
}

func main() {
	options := parseCliArgs()

	newTemplate, err := newStackTemplateFile(options.SourceName, options.Template)
	if err != nil {
		fmt.Printf("Erroring getting the template for cloning. Error: %v", err)
		os.Exit(1)
	}
	defer os.Remove(newTemplate)

	parameters, err := stackParameters(options.SourceName)
	if err != nil {
		fmt.Printf("Error getting source stack parameters. Error: %v", err)
		os.Exit(1)
	}

	for k, v := range cliParameters(options.Attributes) {
		parameters[k] = v
	}

	fmt.Printf("Our merged parameters are %v\n", parameters)

	fmt.Println("Going to clone")

	createStack(options.NewName, parameters, newTemplate)

	os.Exit(0)
}
