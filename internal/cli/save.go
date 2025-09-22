package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/slot"
	"github.com/idelchi/slot/internal/store"
)

// Save returns the cobra command for saving command slots.
func Save(config *string) *cobra.Command {
	var (
		tags        []string
		description string
		force       bool
	)

	cmd := &cobra.Command{
		Use:   "save <slot> <command>",
		Short: "Save a slot",
		Long: heredoc.Doc(`
			Save a command template with optional tags for later execution.

			Commands can include Go template variables like {{.file}} or {{.env}} that will be
			replaced with values when rendering.
		`),
		//nolint:dupword	// False warning
		Example: heredoc.Doc(`
			# Save a simple command
			slot save hello 'echo "Hello World!"'

			# Save with template variables and tags
			slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod

			# Overwrite existing slot
			slot save deploy 'kubectl apply -f {{.file}} --namespace {{.ns}}' --force

			# Save a slot that outputs the content of the slots file
			slot save slots 'cat $(slot ls | tail -1)'
		`),
		Args: cobra.ExactArgs(2), //nolint:mnd   // Clear from the context
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := store.New(*config)
			if err != nil {
				return err
			}

			slots, err := store.Load()
			if err != nil {
				return err
			}

			name, rawCommand := args[0], args[1]
			if slots.Exists(name) && !force {
				return fmt.Errorf("slot %q exists (use --force)", name)
			}

			slots.Delete(name)

			slots.Add(slot.Slot{
				Name:        name,
				Description: description,
				Cmd:         rawCommand,
				Tags:        tags,
			})

			if err := store.Save(slots); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "saved %q\n", name)

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags for the slot (repeatable)")
	cmd.Flags().StringVar(&description, "description", "", "description for the slot")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing slot")

	return cmd
}
