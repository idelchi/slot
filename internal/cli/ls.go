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

			var items []model.Slot
			for _, slot := range database.Slots {
				items = append(items, slot)
			}

			items = filterSlotsByTags(items, filterTags)

			sort.Slice(items, func(i, j int) bool {
				return items[i].Name < items[j].Name
			})

			writer := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, TabSpacing, ' ', 0)
			if _, err := fmt.Fprintln(writer, "NAME\tTAGS\tCMD"); err != nil {
				return err
			}

			for _, slot := range items {
				fmt.Fprintf(writer, "%s\t%s\t%s\n", slot.Name, strings.Join(slot.Tags, ","), slot.Cmd)
			}

			fmt.Fprintf(writer, "\n%s\n", filepath.ToSlash(store.Path))

			return writer.Flush()
		},
	}

	cmd.Flags().StringSliceVar(&filterTags, "tag", nil, "filter by tag (repeatable)")

	return cmd
}

// filterSlotsByTags returns slots that contain all specified tags.
func filterSlotsByTags(slots []model.Slot, filterTags []string) []model.Slot {
	if len(filterTags) == 0 {
		return slots
	}

	var result []model.Slot

	for _, slot := range slots {
		hasAllTags := true

		for _, filterTag := range filterTags {
			found := false

			for _, slotTag := range slot.Tags {
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
			result = append(result, slot)
		}
	}

	return result
}
