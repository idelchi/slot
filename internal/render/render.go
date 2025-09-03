// Package render provides Go template parsing and execution for command string substitution.
package render

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

// Apply executes a Go template with provided variables, returning an error if parsing fails or variables are missing.
func Apply(templateString string, variables map[string]string) (string, error) {
	template, err := template.New("cmd").Funcs(sprig.FuncMap()).Option("missingkey=error").Parse(templateString)
	if err != nil {
		return "", err
	}

	// Execute the template with variables
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, variables); err != nil {
		return "", errToMissingKey(err)
	}

	return strings.TrimSpace(buffer.String()), nil
}

// errToMissingKey formats the original error from text/template to a friendler one.
func errToMissingKey(err error) error {
	message := err.Error()

	seek := "map has no entry for key "

	if strings.Contains(message, seek) {
		_, variable, _ := strings.Cut(message, seek)

		return fmt.Errorf("missing template variable: %s", variable)
	}

	return err
}
