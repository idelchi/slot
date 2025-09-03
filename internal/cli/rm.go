package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/store"
)

// Remove returns the cobra command for removing command slots.
func Remove() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rm <slot>",
		Short:   "Delete a slot",
		Aliases: []string{"remove", "delete"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := store.New()
			if err != nil {
				return err
			}

			database, err := store.Load()
			if err != nil {
				return err
			}

			slot := args[0]
			if _, ok := database.Slots[slot]; !ok {
				return fmt.Errorf("no such slot %q", slot)
			}

			delete(database.Slots, slot)

			if err := store.Save(database); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed %q\n", slot)

			return nil
		},
	}

	return cmd
}
