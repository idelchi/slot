package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/store"
)

// Remove returns the cobra command for removing command slots.
func Remove(config *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "remove <slot>",
		Short:   "Delete a slot",
		Aliases: []string{"rm"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := store.New(*config)
			if err != nil {
				return err
			}

			slots, err := store.Load()
			if err != nil {
				return err
			}

			if len(slots) == 0 {
				return errors.New("no slots to remove")
			}

			slot := args[0]
			if !slots.Exists(slot) {
				return fmt.Errorf("no such slot %q: did you mean %q?", slot, slots.Closest(slot))
			}

			slots.Delete(slot)

			if err := store.Save(slots); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "removed %q\n", slot)

			return nil
		},
	}

	return cmd
}
