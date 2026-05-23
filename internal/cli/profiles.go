package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/SciTee/ledgerly/internal/profile"
	"github.com/spf13/cobra"
)

func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Inspect the bank profiles ledgerly knows about",
	}
	cmd.AddCommand(newProfilesListCmd())
	return cmd
}

func newProfilesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List the built-in bank profiles",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			profiles, err := profile.Builtins()
			if err != nil {
				return err
			}
			tw := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 2, 2, ' ', 0)
			fmt.Fprintln(tw, "NAME\tCURRENCY\tDELIMITER\tDECIMAL")
			for _, p := range profiles {
				delim := p.Delimiter
				if delim == "" {
					delim = "auto"
				}
				dec := p.Decimal
				if dec == "" {
					dec = "auto"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", p.Name, p.Currency, delim, dec)
			}
			return tw.Flush()
		},
	}
}
