package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/render"
	"github.com/idelchi/slot/internal/store"
)

// Render returns the cobra command for rendering command slots.
func Render() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "render <slot> [key=value...]",
		Short: "Render a slot",
		Long: heredoc.Doc(`
			Render a saved command slot, substituting template variables with provided values.

			Templates use Go template syntax: {{.variable}} is replaced with values from key=value arguments.

			The rendered command is printed to stdout for shell integration.
		`),
		Example: heredoc.Doc(`
			# Render a command with variable substitution
			slot render deploy file=k8s.yml ns=production

			# Render command without variables
			slot render hello
		`),
		Args: cobra.MinimumNArgs(1),
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
			if !database.Exists(slot) {
				return fmt.Errorf("no such slot %q", slot)
			}

			withs := args[1:]
			variables, err := parseWiths(withs)
			if err != nil {
				return err
			}

			rendered, err := render.Apply(database.Get(slot).Cmd, variables)
			if err != nil {
				return err
			}

			if _, err := fmt.Fprintln(cmd.OutOrStdout(), rendered); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

// parseWiths parses key=value pairs into a key-value map.
func parseWiths(keyValues []string) (map[string]string, error) {
	var errs []error

	out := make(map[string]string)

	for _, keyValue := range keyValues {
		key, value, found := strings.Cut(keyValue, "=")

		if !found {
			errs = append(errs, fmt.Errorf("missing value: %q", keyValue))

			continue
		}

		if key == "" {
			errs = append(errs, fmt.Errorf("missing key: %q", keyValue))

			continue
		}

		out[key] = value
	}

	return out, errors.Join(errs...)
}
