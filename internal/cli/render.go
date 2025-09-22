package cli

import (
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/render"
	"github.com/idelchi/slot/internal/store"
)

// Render returns the cobra command for rendering command slots.
func Render(config *string) *cobra.Command {
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
			args, afterDash := splitAtDash(cmd, args)
			if len(args) < 1 {
				return errors.New("requires at least 1 arg(s), only received 0")
			}

			store, err := store.New(*config)
			if err != nil {
				return err
			}

			slots, err := store.Load()
			if err != nil {
				return err
			}

			if len(slots) == 0 {
				return errors.New("no slots to render")
			}

			slot := args[0]
			if !slots.Exists(slot) {
				return fmt.Errorf("no such slot %q: did you mean %q?", slot, slots.Closest(slot))
			}

			variables := map[string]any{}

			variables["SLOTS_FILE"] = filepath.ToSlash(store.Path())
			variables["SLOTS_DIR"] = filepath.ToSlash(filepath.Dir(store.Path()))
			variables["CLI_ARGS"] = strings.Join(afterDash, " ")

			withs, err := parseWiths(args[1:])
			if err != nil {
				return err
			}

			maps.Copy(variables, withs)

			//nolint:errcheck,forcetypeassert  // args are always strings
			variables["CLI_ARGS_SPLIT"] = strings.Split(variables["CLI_ARGS"].(string), " ")

			rendered, err := render.Apply(slots.Get(slot).Cmd, variables)
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
func parseWiths(keyValues []string) (map[string]any, error) {
	var errs []error

	out := make(map[string]any)

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

// splitAtDash splits args at the first occurrence of "--".
func splitAtDash(cmd *cobra.Command, args []string) (beforeDash, afterDash []string) {
	n := cmd.ArgsLenAtDash()
	if n >= 0 {
		return args[:n], args[n:]
	}

	return args, nil
}
