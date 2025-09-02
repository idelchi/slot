package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/store"
)

// Constants for CLI configuration.
const (
	KeyValueParts = 2
	TabSpacing    = 2
	SaveArgsCount = 2
)

// Options represents the root level configuration for the CLI application.
type Options struct {
	// Verbose enables verbose output.
	Verbose bool
}

// Execute runs the root command for the slot CLI application.
func Execute(version string) error {
	options := &Options{}

	root := &cobra.Command{
		Use:   "slot",
		Short: "Save and render named shell command slots",
		Long: heredoc.Doc(`
			Slot is a CLI tool for managing named shell command slots with Go template substitution.

			Save commands with Go template variables and tags, then render them with variable substitution.
			Use shell completions to place rendered commands into your prompt for execution.
			All operations are logged for audit purposes.
		`),
		Example: heredoc.Doc(`
			# Save a command with template variables
			$ slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod

			# Render with variable substitution
			$ slot run deploy --with file=k8s.yml

			# Generate shell integration
			$ slot completions bash

			# List all slots
			$ slot ls
		`),
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.SetVersionTemplate("{{ .Version }}\n")
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.Flags().SortFlags = false
	root.PersistentFlags().SortFlags = false

	root.CompletionOptions.DisableDefaultCmd = true
	cobra.EnableCommandSorting = false

	root.PersistentFlags().
		BoolVarP(&options.Verbose, "verbose", "v", false, "Increase verbosity level")

	root.AddCommand(
		Save(options),
		Run(options),
		List(options),
		Remove(options),
		Path(options),
		Completions(options),
	)

	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}

// mustStore creates a new Store instance or exits on error.
func mustStore() *store.Store {
	s, err := store.New()
	if err != nil {
		panic(err) // In CLI context, this is acceptable
	}

	return s
}

// parseWith parses --with flags into a key-value map.
func parseWith(items []string) (map[string]string, error) {
	out := make(map[string]string)

	for _, s := range items {
		kv := strings.SplitN(s, "=", KeyValueParts)
		if len(kv) != KeyValueParts {
			return nil, fmt.Errorf("bad --with format %q (want KEY=VAL)", s)
		}

		out[kv[0]] = kv[1]
	}

	return out, nil
}
