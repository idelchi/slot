package slot

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// Common header for both outputs.
const slotsHeader = "NAME\tCMD\tTAGS\tDESCRIPTION"

// makeRecords builds rows with configurable newline handling and trimming.
func makeRecords(slots Slots, writer io.Writer) error {
	const newlineReplacement = "^J"

	records := make([][]string, 0, len(slots))

	for _, slot := range slots {
		cmd := strings.ReplaceAll(slot.Cmd, "\n", newlineReplacement)

		records = append(records, []string{
			slot.Name,
			cmd,
			strings.Join(slot.Tags, ","),
			slot.Description,
		})
	}

	for _, record := range records {
		if _, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", record[0], record[1], record[2], record[3]); err != nil {
			return err
		}
	}

	return nil
}

// asTable writes human-readable aligned columns via tabwriter.
func asTable(slots Slots, writer io.Writer) error {
	const TabSpacing = 2

	tabWriter := tabwriter.NewWriter(writer, 0, 0, TabSpacing, ' ', 0)

	if _, err := fmt.Fprintln(tabWriter, slotsHeader); err != nil {
		return err
	}

	if err := makeRecords(slots, tabWriter); err != nil {
		return err
	}

	return tabWriter.Flush()
}

// asTSV writes deterministic machine-readable TSV.
func asTSV(slots Slots, writer io.Writer) error {
	if _, err := fmt.Fprintln(writer, slotsHeader); err != nil {
		return err
	}

	return makeRecords(slots, writer)
}
