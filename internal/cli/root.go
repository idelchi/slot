package cli

import (
	"os"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/store"
)

// Execute runs the root command for the slot CLI application.
func Execute(version string) error {
	root := &cobra.Command{
		Use:   "slot",
		Short: "Manage named shell command slots",
		Long: heredoc.Doc(`
			Slot is a CLI tool for managing named shell command slots with Go template substitution.

			Save commands with Go template variables and tags, then render them with variable substitution.
			Use shell integration to place rendered commands into your prompt for execution.

			Add 'eval "$(slot init <shell>)"' to your shell configuration to enable command substitution
			with 'slot run <slot>'.

			"slot init <shell> --fzf" allows for further key-bindings to Ctrl-Z and Ctrl-X.
		`),
		//nolint:dupword	// False warning
		Example: heredoc.Doc(`
			# Save a command with template variables
			slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod

			# List all slots
			slot list

			# Render with variable substitution
			slot render deploy file=k8s.yml

			# Remove a slot
			slot remove deploy

			# Generate shell integration
			slot init zsh
		`),
		Version:       version,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	root.SetVersionTemplate("{{ .Version }}\n")
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	root.Flags().SortFlags = false
	root.PersistentFlags().SortFlags = false

	cobra.EnableCommandSorting = false

	config := os.Getenv("SLOTS_FILE")
	if config == "" {
		config, _ = store.DefaultSlotsFile()
	}

	root.PersistentFlags().StringVar(&config, "config", config, "path to the configuration file")

	root.AddCommand(
		Save(&config),
		Render(&config),
		List(&config),
		Remove(&config),
		Path(&config),
		Init(),
	)

	return root.Execute()
}
