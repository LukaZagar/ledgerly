// Package cli wires the ledgerly subcommands together on top of cobra.
package cli

import "github.com/spf13/cobra"

// NewRootCmd builds the root ledgerly command tree. version is shown by
// --version and is injected from main at build time.
func NewRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "ledgerly",
		Short: "Unify and categorize bank CSV exports",
		Long: "Ledgerly converts the wildly different CSV exports of various banks into\n" +
			"one clean, unified and automatically categorized transaction list.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	root.SetVersionTemplate("ledgerly {{.Version}}\n")

	root.AddCommand(newConvertCmd())
	root.AddCommand(newProfilesCmd())
	root.AddCommand(newSummaryCmd())
	return root
}

// Execute builds the command tree and runs it, returning any error so main can
// set the process exit code.
func Execute(version string) error {
	return NewRootCmd(version).Execute()
}
