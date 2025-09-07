package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/slots"
	"github.com/idelchi/slot/internal/store"
)

// Save returns the cobra command for saving command slots.
func Save() *cobra.Command {
	var (
		tags  []string
		force bool
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
		Args: cobra.ExactArgs(SaveArgsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := store.New()
			if err != nil {
				return err
			}

			database, err := store.Load()
			if err != nil {
				return err
			}

			slot, rawCommand := args[0], args[1]
			if database.Exists(slot) && !force {
				return fmt.Errorf("slot %q exists (use --force)", slot)
			}

			database.Delete(slot)

			database.Add(slots.Slot{
				Name: slot,
				Cmd:  rawCommand,
				Tags: tags,
			})

			if err := store.Save(database); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "saved %q\n", slot)

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags for the slot (repeatable)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing slot")

	return cmd
}
