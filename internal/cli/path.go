package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Path returns the cobra command for showing the slots file path.
func Path(_ *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Show the path to the slots file",
		Example: heredoc.Doc(`
			# Show slots file path
			$ slot path
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			store := mustStore()
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), store.Path); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
