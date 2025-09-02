package cli

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/completions"
)

// Completions returns the cobra command for generating shell completions.
func Completions(_ *Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completions <shell>",
		Short: "Generate shell completion snippets",
		Long: heredoc.Doc(`
			Generate shell integration snippets for the specified shell.

			The integration allows 'slot run' to place rendered commands into
			the prompt for editing before execution.
		`),
		Example: heredoc.Doc(`
			# Generate bash completion
			$ slot completions bash >> ~/.bashrc

			# Generate zsh completion
			$ slot completions zsh >> ~/.zshrc
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			shell := args[0]

			switch shell {
			case "bash":
				fmt.Fprint(cmd.OutOrStdout(), completions.Bash)

				return nil
			case "zsh":
				fmt.Fprint(cmd.OutOrStdout(), completions.Zsh)

				return nil
			default:
				return fmt.Errorf("unsupported shell %q (supported: bash, zsh)", shell)
			}
		},
	}

	return cmd
}
