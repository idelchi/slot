// Package render provides Go template parsing and execution for command string substitution.
package render

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// templateRegex matches Go template variable patterns in the format {{.Variable}}.
var templateRegex = regexp.MustCompile(`\{\{\s*\.([A-Za-z_][A-Za-z0-9_]*)\s*\}\}`)

// Placeholders extracts all unique placeholder names from a template string.
func Placeholders(template string) ([]string, error) {
	matches := templateRegex.FindAllStringSubmatch(template, -1)
	seen := make(map[string]bool)

	var names []string

	for _, match := range matches {
		if len(match) > 1 {
			name := match[1]
			if !seen[name] {
				seen[name] = true
				names = append(names, name)
			}
		}
	}

	sort.Strings(names)

	return names, nil
}

// Apply executes a Go template with provided variables, returning an error if parsing fails or variables are missing.
func Apply(templateString string, variables map[string]string) (string, error) {
	// Check if all required variables are provided
	placeholders, err := Placeholders(templateString)
	if err != nil {
		return "", err
	}

	var missing []string

	for _, placeholder := range placeholders {
		if _, ok := variables[placeholder]; !ok {
			missing = append(missing, placeholder)
		}
	}

	if len(missing) > 0 {
		return "", fmt.Errorf("missing template variables: %s", //nolint:err113 // Dynamic error is fine
			strings.Join(missing, ", "))
	}

	// Create and parse the template
	template, err := template.New("cmd").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	// Execute the template with variables
	var buffer bytes.Buffer
	if err := template.Execute(&buffer, variables); err != nil {
		return "", fmt.Errorf("template execution error: %w", err)
	}

	return strings.TrimSpace(buffer.String()), nil
}
