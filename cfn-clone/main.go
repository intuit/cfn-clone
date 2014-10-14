package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
)

func prettyParameters(params map[string]string) string {
	var b bytes.Buffer
	if len(params) > 0 {
		w := new(tabwriter.Writer)
		w.Init(&b, 0, 8, 0, '\t', 0)

		keys := []string{}
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.Write([]byte("The merged parameters are:\n"))
		for _, k := range keys {
			fmt.Fprintf(w, "%s \t%s\n", k, params[k])
		}
		w.Flush()
	}
	return b.String()
}

func main() {
	options := parseCliArgs()

	newTemplate, err := newStackTemplateFile(options.SourceName, options.Template)
	if err != nil {
		fmt.Printf("Erroring getting the template for cloning. %s\n", err.Error())
		os.Exit(1)
	}
	defer os.Remove(newTemplate)

	parameters, err := stackParameters(options.SourceName)
	if err != nil {
		fmt.Printf("Error getting source stack parameters. %s\n", err.Error())
		os.Exit(1)
	}

	for k, v := range paramsFromCli(options.Attributes) {
		parameters[k] = v
	}

	fmt.Println(prettyParameters(parameters))

	fmt.Println("Going to clone")

	output, err := createStack(options.NewName, parameters, newTemplate)
	if err != nil {
		fmt.Printf("Unable to create new stack. %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Success with output '%s'.\n", output)
	os.Exit(0)
}
