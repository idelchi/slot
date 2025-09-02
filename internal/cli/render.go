package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/render"
)

// Run returns the cobra command for rendering command slots.
func Run(_ *Options) *cobra.Command {
	var withs []string

	cmd := &cobra.Command{
		Use:   "render <name>",
		Short: "Render a saved command slot",
		Long: heredoc.Doc(`
			Render a saved command slot, substituting template variables with provided values.

			Templates use Go template syntax: {{.variable}} is replaced with values from --with flags.
			The rendered command is printed to stdout for shell integration.
		`),
		Example: heredoc.Doc(`
			# Render a command with variable substitution
			$ slot render deploy --with file=k8s.yml --with ns=production

			# Render command without variables
			$ slot render hello
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustStore()
			name := args[0]

			database, err := store.Load()
			if err != nil {
				return err
			}

			slot, ok := database.Slots[name]
			if !ok {
				return fmt.Errorf("no such slot %q", name)
			}

			variables, err := parseWith(withs)
			if err != nil {
				return err
			}

			rendered, err := render.Apply(slot.Cmd, variables)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), rendered); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringSliceVar(&withs, "with", nil, "KEY=VAL pairs for placeholder substitution (repeatable)")

	return cmd
}
