package cli

import (
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

			if tsv {
				return slots.Render("tsv", cmd.OutOrStdout())
			}

			return slots.Render("table", cmd.OutOrStdout())
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
