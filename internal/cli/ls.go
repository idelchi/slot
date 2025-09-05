package cli

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/idelchi/slot/internal/model"
	"github.com/idelchi/slot/internal/store"
)

// List returns the cobra command for listing command slots.
func List() *cobra.Command {
	var filterTags []string

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List saved slots",
		Long: heredoc.Doc(`
			List all saved command slots with their names, tags, and commands.
		`),
		Example: heredoc.Doc(`
			# List all slots in table format
			slot ls

			# Show only slots tagged with 'k8s'
			slot ls --tag k8s

			# Multiple tag filters (AND logic)
			slot ls --tag k8s --tag prod
		`),
		Aliases: []string{"list"},
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

			var items []slotWithName
			for name, slot := range database.Slots {
				items = append(items, slotWithName{name: name, slot: slot})
			}

			items = filterSlotsByTags(items, filterTags)

			sort.Slice(items, func(i, j int) bool {
				return items[i].name < items[j].name
			})

			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, TabSpacing, ' ', 0)
			if _, err := fmt.Fprintln(writer, "NAME\tTAGS\tCMD"); err != nil {
				return err
			}

			for _, item := range items {
				// Replace newlines with ^J (caret notation for newline) for display
				displayCmd := strings.ReplaceAll(item.slot.Cmd, "\n", "^J")
				displayCmd = strings.TrimSpace(displayCmd)
				fmt.Fprintf(writer, "%s\t%s\t%s\n", item.name, strings.Join(item.slot.Tags, ","), displayCmd)
			}

			fmt.Fprintf(writer, "\n%s\n", filepath.ToSlash(store.Path))

			return writer.Flush()
		},
	}

	cmd.Flags().StringSliceVar(&filterTags, "tag", nil, "filter by tag (repeatable)")

	return cmd
}

// slotWithName pairs a slot name with its data for sorting and display.
type slotWithName struct {
	name string
	slot model.Slot
}

// filterSlotsByTags returns slots that contain all specified tags.
func filterSlotsByTags(slots []slotWithName, filterTags []string) []slotWithName {
	if len(filterTags) == 0 {
		return slots
	}

	var result []slotWithName

	for _, item := range slots {
		hasAllTags := true

		for _, filterTag := range filterTags {
			found := false

			for _, slotTag := range item.slot.Tags {
				if slotTag == filterTag {
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
			result = append(result, item)
		}
	}

	return result
}
