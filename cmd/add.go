package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <name> <oci-url>",
		Short: "Bookmark an OCI chart reference",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, url := args[0], args[1]
			if err := getStore().Add(name, url); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added %q -> %s\n", name, url)
			return nil
		},
	}
}
