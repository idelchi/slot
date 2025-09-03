package cli

import (
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
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
	root := &cobra.Command{
		Use:   "slot",
		Short: "Manage named shell command slots",
		Long: heredoc.Doc(`
			Slot is a CLI tool for managing named shell command slots with Go template substitution.

			Save commands with Go template variables and tags, then render them with variable substitution.
			Use shell integration to place rendered commands into your prompt for execution.

			Add 'eval $(slot init <shell>)' to your shell configuration to enable command substitution
			with 'slot run <slot>'.
		`),
		Example: heredoc.Doc(`
			# Save a command with template variables
			slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod

			# Render with variable substitution
			slot render deploy file=k8s.yml

			# Generate shell integration
			slot init zsh

			# List all slots
			slot ls
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

	root.AddCommand(
		Save(),
		Render(),
		List(),
		Remove(),
		Init(),
	)

	return root.Execute()
}
