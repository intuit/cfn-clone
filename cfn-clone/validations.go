package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

func validateCliExists(cmd string) error {
	_, err := exec.LookPath(cmd)
	return err
}

func validateCliParameters(params []string) error {
	for _, p := range params {
		v := strings.SplitN(p, "=", 2)
		if len(v) != 2 {
			return errors.New("Attribute '" + p + "' must be '=' separated key, value")
		}
	}
	return nil
}

func validateTemplateExists(path string) error {
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func validateSourceStackExists(name string) error {
	args := []string{
		"cloudformation",
		"describe-stacks",
		"--stack-name",
		name,
	}

	cmd := exec.Command("aws", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("Error verifying source stack. Error: " + string(output))
	}

	return nil
}
