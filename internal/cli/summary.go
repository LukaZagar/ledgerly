package cli

import (
	"github.com/SciTee/ledgerly/internal/summary"
	"github.com/spf13/cobra"
)

func newSummaryCmd() *cobra.Command {
	var (
		profileName string
		autoDetect  bool
		account     string
		rulesFile   string
		noCat       bool
	)

	cmd := &cobra.Command{
		Use:   "summary [flags] <file.csv> [more.csv...]",
		Short: "Print category and monthly totals without writing any file",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := selectProfile(profileName, autoDetect, args[0])
			if err != nil {
				return err
			}
			res, err := runPipeline(args, pipelineOptions{
				profile:      p,
				account:      account,
				rulesFile:    rulesFile,
				noCategorize: noCat,
			})
			if err != nil {
				return err
			}
			return summary.Compute(res.Transactions).WriteText(cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringVarP(&profileName, "profile", "p", "", "bank profile: a built-in name or a path to a YAML file")
	cmd.Flags().BoolVar(&autoDetect, "auto-detect", false, "auto-detect the profile from the CSV header")
	cmd.Flags().StringVarP(&account, "account", "a", "", "account label for the rows (default: input file name)")
	cmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "user categorization rules file (stacked on top of defaults)")
	cmd.Flags().BoolVar(&noCat, "no-categorize", false, "skip categorization")

	return cmd
}
