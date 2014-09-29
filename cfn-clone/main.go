package main

import (
	"fmt"
	"os"
)

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

	for k, v := range paramsFromCli(options.Attributes) {
		parameters[k] = v
	}

	fmt.Printf("Our merged parameters are %v\n", parameters)

	fmt.Println("Going to clone")

	output, err := createStack(options.NewName, parameters, newTemplate)
	if err != nil {
		fmt.Printf("Received error '%s' with output '%s'.\n", err.Error(), output)
		os.Exit(1)
	}

	fmt.Printf("Success with output '%s'.\n", output)
	os.Exit(0)
}
