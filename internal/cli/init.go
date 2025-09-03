package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/completions"
)

// Init returns the cobra command for generating shell integration scripts.
func Init() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init <shell>",
		Short: "Generate shell integration snippets",
		Long: heredoc.Doc(`
			Generate shell integration snippets for the specified shell.

			The integration allows 'slot run' to place rendered commands into
			the prompt for editing before execution.
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			supported := []string{"bash", "zsh"}

			shell := args[0]

			switch shell {
			case "bash":
				fmt.Fprint(cmd.OutOrStdout(), completions.Bash)

				return nil
			case "zsh":
				fmt.Fprint(cmd.OutOrStdout(), completions.Zsh)

				return nil
			default:
				return fmt.Errorf("unsupported shell %q (supported: %v)", shell, strings.Join(supported, ", "))
			}
		},
	}

	return cmd
}
