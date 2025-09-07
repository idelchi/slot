package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/slots"
	"github.com/idelchi/slot/internal/store"
)

// List returns the cobra command for listing command slots.
func List() *cobra.Command {
	var (
		filterTags []string
		tsv        bool
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List saved slots",
		Long: heredoc.Doc(`
			List all saved command slots with their names, tags, and commands.
		`),
		Example: heredoc.Doc(`
			# List all slots in table format
			slot list

			# Show only slots tagged with 'k8s'
			slot list --tag k8s

			# Multiple tag filters (AND logic)
			slot list --tag k8s --tag prod
		`),
		Aliases: []string{"ls"},
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			store, err := store.New()
			if err != nil {
				return err
			}

			database, err := store.Load()
			if err != nil {
				return err
			}

			slots := filterSlotsByTags(database, filterTags)

			// Machine-readable output for fzf/zsh parsing
			if tsv {
				if _, err := fmt.Fprintln(cmd.OutOrStdout(), "NAME\tTAGS\tCMD"); err != nil {
					return err
				}

				// Deterministic TSV: name \t tags(csv) \t cmd(with \n escaped) \n
				for _, slot := range slots {
					cmdStr := strings.ReplaceAll(slot.Cmd, "\n", `\n`)
					// No header, no footer/path line
					_, err := fmt.Fprintf(
						cmd.OutOrStdout(),
						"%s\t%s\t%s\n",
						slot.Name,
						strings.Join(slot.Tags, ","),
						cmdStr,
					)
					if err != nil {
						return err
					}
				}

				return nil
			}

			// Human-readable table (default)
			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, TabSpacing, ' ', 0)
			if _, err := fmt.Fprintln(writer, "NAME\tTAGS\tCMD"); err != nil {
				return err
			}

			for _, slot := range slots {
				// Replace newlines with ^J (caret notation for newline) for display
				displayCmd := strings.ReplaceAll(slot.Cmd, "\n", "^J")
				displayCmd = strings.TrimSpace(displayCmd)
				fmt.Fprintf(writer, "%s\t%s\t%s\n", slot.Name, strings.Join(slot.Tags, ","), displayCmd)
			}

			fmt.Fprintf(writer, "\n%s\n", filepath.ToSlash(store.Path))

			return writer.Flush()
		},
	}

	cmd.Flags().StringSliceVar(&filterTags, "tag", nil, "filter by tag (repeatable)")
	cmd.Flags().BoolVar(&tsv, "tsv", false, "output in TSV format")

	return cmd
}

// filterSlotsByTags returns slots that contain all specified tags.
func filterSlotsByTags(database slots.Slots, filterTags []string) slots.Slots {
	if len(filterTags) == 0 {
		return database
	}

	var result slots.Slots

	for _, slot := range database {
		hasAllTags := true

		for _, filterTag := range filterTags {
			found := false

			for _, tag := range slot.Tags {
				if tag == filterTag {
					found = true

					break
				}
			}

			if !found {
				hasAllTags = false

				break
			}
		}

		if hasAllTags {
			result = append(result, slot)
		}
	}

	return result
}
