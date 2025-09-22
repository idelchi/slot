package cli

import (
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
)

// Path returns the cobra command for displaying the path to the slot store.
func Path(slotsFile *string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "path",
		Short: "Display the path to the slot store",
		Long: heredoc.Doc(`
			Display the absolute path to the directory where slot files are stored.
		`),
		Example: heredoc.Doc(`
			# Show the path to the slot store
			slot path
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), filepath.FromSlash(*slotsFile)); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
