package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/SciTee/ledgerly/internal/export"
	"github.com/SciTee/ledgerly/internal/model"
	"github.com/SciTee/ledgerly/internal/parse"
	"github.com/SciTee/ledgerly/internal/profile"
	"github.com/spf13/cobra"
)

func newConvertCmd() *cobra.Command {
	var (
		profileName string
		autoDetect  bool
		account     string
		format      string
		outPath     string
		rulesFile   string
		noCat       bool
	)

	cmd := &cobra.Command{
		Use:   "convert [flags] <file.csv> [more.csv...]",
		Short: "Convert one or more bank CSVs into the unified format",
		Long: "Reads the given bank CSV files using the selected profile, merges them\n" +
			"into a single chronological timeline, removes duplicate rows and writes\n" +
			"the unified result as CSV or JSON.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if format != "csv" && format != "json" {
				return fmt.Errorf("unknown format %q (use csv or json)", format)
			}

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

			out, closeOut, err := openOutput(outPath)
			if err != nil {
				return err
			}
			defer closeOut()

			if err := write(out, format, res.Transactions); err != nil {
				return err
			}

			fmt.Fprintf(cmd.ErrOrStderr(), "converted %d transactions", len(res.Transactions))
			if res.Duplicates > 0 {
				fmt.Fprintf(cmd.ErrOrStderr(), " (%d duplicates removed)", res.Duplicates)
			}
			fmt.Fprintln(cmd.ErrOrStderr())
			return nil
		},
	}

	cmd.Flags().StringVarP(&profileName, "profile", "p", "", "bank profile: a built-in name or a path to a YAML file")
	cmd.Flags().BoolVar(&autoDetect, "auto-detect", false, "auto-detect the profile from the CSV header")
	cmd.Flags().StringVarP(&account, "account", "a", "", "account label for the rows (default: input file name)")
	cmd.Flags().StringVarP(&format, "format", "f", "csv", "output format: csv or json")
	cmd.Flags().StringVarP(&outPath, "out", "o", "", "output file (default: stdout)")
	cmd.Flags().StringVarP(&rulesFile, "rules", "r", "", "user categorization rules file (stacked on top of defaults)")
	cmd.Flags().BoolVar(&noCat, "no-categorize", false, "skip categorization")

	return cmd
}

// selectProfile resolves which profile to use: an explicit --profile, or
// --auto-detect which reads the header of the first input file and matches it
// against the built-in signatures.
func selectProfile(name string, autoDetect bool, firstInput string) (*profile.Profile, error) {
	switch {
	case autoDetect:
		f, err := os.Open(firstInput)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		tab, err := parse.ReadCSV(f)
		if err != nil {
			return nil, err
		}
		return profile.Detect(tab.Header)
	case name != "":
		return profile.Resolve(name)
	default:
		return nil, fmt.Errorf("a profile is required: pass --profile or --auto-detect")
	}
}

func write(w io.Writer, format string, txs []model.Transaction) error {
	if format == "json" {
		return export.WriteJSON(w, txs)
	}
	return export.WriteCSV(w, txs)
}

// openOutput returns the writer for the chosen output path. When the path is
// empty it writes to stdout and the returned close func is a no-op.
func openOutput(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, func() {}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, err
	}
	return f, func() { f.Close() }, nil
}
