package cli

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/integration"
)

// Init returns the cobra command for generating shell integration scripts.
func Init() *cobra.Command {
	var fzf bool

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

			inits := []string{}

			switch shell {
			case "bash":
				inits = append(inits, integration.Bash)

				if fzf {
					inits = append(inits, integration.BashFzf)
				}
			case "zsh":
				inits = append(inits, integration.Zsh)

				if fzf {
					inits = append(inits, integration.ZshFzf)
				}
			default:
				return fmt.Errorf("unsupported shell %q (supported: %v)", shell, strings.Join(supported, ", "))
			}

			_, err := fmt.Fprintln(cmd.OutOrStdout(), strings.Join(inits, "\n"))

			return err
		},
	}

	cmd.Flags().BoolVar(&fzf, "fzf", false, "enable fzf support and bind to Ctrl-X and Ctrl-Z keys (zsh only)")

	return cmd
}
