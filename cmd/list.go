package cmd

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all bookmarked OCI chart references",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			items, err := getStore().List()
			if err != nil {
				return err
			}

			if len(items) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No bookmarks found. Use 'helm oci add <name> <oci-url>' to add one.")
				return nil
			}

			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "NAME\tURL")
			for _, b := range items {
				fmt.Fprintf(w, "%s\t%s\n", b.Name, b.URL)
			}
			return w.Flush()
		},
	}
}
