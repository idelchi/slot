package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/model"
)

// Save returns the cobra command for saving command slots.
func Save(_ *Options) *cobra.Command {
	var (
		tags  []string
		force bool
	)

	cmd := &cobra.Command{
		Use:   "save <name> <command>",
		Short: "Save a named command slot",
		Long: heredoc.Doc(`
			Save a command template with optional tags for later execution.

			Commands can include Go template variables like {{.file}} or {{.env}} that will be
			replaced with values when running the slot.
		`),
		Example: heredoc.Doc(`
			# Save a simple command
			$ slot save hello 'echo "Hello World!"'

			# Save with template variables and tags
			$ slot save deploy 'kubectl apply -f {{.file}}' --tags k8s --tags prod

			# Overwrite existing slot
			$ slot save deploy 'kubectl apply -f {{.file}} --namespace {{.ns}}' --force
		`),
		Args: cobra.ExactArgs(SaveArgsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustStore()
			name, rawCommand := args[0], args[1]

			// Normalize line endings: CRLF/CR -> LF
			// rawCommand = strings.ReplaceAll(rawCommand, "\r\n", "\n")
			// rawCommand = strings.ReplaceAll(rawCommand, "\r", "\n")

			database, err := store.Load()
			if err != nil {
				return err
			}

			if database.Slots == nil {
				database.Slots = make(map[string]model.Slot)
			}

			if _, exists := database.Slots[name]; exists && !force {
				return fmt.Errorf("slot %q exists (use --force)", name)
			}

			slot := model.Slot{
				Name: name,
				Cmd:  rawCommand,
				Tags: tags,
			}

			database.Slots[name] = slot

			if err := store.Save(database); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "saved %s\n", name)

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&tags, "tags", nil, "tags for the slot (repeatable)")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing slot")

	return cmd
}
