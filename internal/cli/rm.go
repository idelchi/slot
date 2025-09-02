package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Remove returns the cobra command for removing command slots.
func Remove(_ *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm <name>",
		Short:   "Delete a saved slot",
		Aliases: []string{"remove", "delete"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := mustStore()
			name := args[0]

			database, err := store.Load()
			if err != nil {
				return err
			}

			if _, ok := database.Slots[name]; !ok {
				return fmt.Errorf("no such slot %q", name)
			}

			delete(database.Slots, name)

			if err := store.Save(database); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed %q\n", name)

			return nil
		},
	}

	return cmd
}
