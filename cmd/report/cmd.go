package report

import (
	"io/ioutil"
	"os"
	"time"

	"github.com/riccardom/briatore/report"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/riccardom/briatore/types"
)

const (
	flagFile   = "file"
	flagOutput = "output"
)

// GetReportCmd returns the command to crete a report for a specific date
func GetReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "report [date]",
		Short:   "Reports the data for the provided addresses",
		Example: "report cosmos-hub osmosis chihuahua",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(os.Stdout)

			cfg, err := types.ReadConfig(cmd)
			if err != nil {
				return err
			}

			date, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return err
			}

			outValue, err := cmd.Flags().GetString(flagOutput)
			if err != nil {
				return err
			}

			out, err := types.ParseOutput(outValue)
			if err != nil {
				return err
			}

			bz, err := report.GetReportBytes(cfg, cfg.Addresses, date, out)
			if err != nil {
				return err
			}

			outputFile, _ := cmd.Flags().GetString(flagFile)
			if outputFile != "" {
				log.Info().Msg("writing reports to file")
				return ioutil.WriteFile(outputFile, bz, 0666)
			}

			cmd.Print(string(bz))

			return nil
		},
	}

	cmd.Flags().String(flagFile, "", "File where to store the reports")
	cmd.Flags().String(flagOutput, types.OutText.String(), "Type of output (supported values: json, text)")

	return cmd
}
