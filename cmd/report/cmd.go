package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/riccardom/briatore/reporter"
	"github.com/riccardom/briatore/types"
)

const (
	flagFile   = "file"
	flagOutput = "output"

	outText = "text"
	outJSON = "json"
)

func GetReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report [year] [[chains]]",
		Short: "Reports the data from the given chains for the given year",
		Long: `
Reports the data from the given chains for the given year.

If no chain is provided, then all chains present inside the configuration file will be reported.
Multiple chains can be specified separating them using spaces.`,
		Example: "report cosmos-hub osmosis chihuahua",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(os.Stdout)

			cfg, err := types.ReadConfig(cmd)
			if err != nil {
				return err
			}

			year, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			chains := args[1:]
			if len(chains) == 0 {
				chains = cfg.GetChainsList()
			}

			if len(chains) == 0 {
				return fmt.Errorf("cannot parse an empty list of chains: check your config and try again")
			}

			log.Debug().Strs("chains", chains).Msg("getting reports")

			cdc, _ := simapp.MakeCodecs()

			reports := make([]*types.ChainReport, len(chains))
			for i, chain := range chains {
				log.Debug().Str("chain", chain).Msg("getting configuration")

				chainCfg, found := cfg.GetChainConfig(chain)
				if !found {
					return fmt.Errorf("config for chain %s not found", chain)
				}

				rep, err := reporter.NewReporter(cfg.Report, chainCfg, cdc)
				if err != nil {
					return err
				}

				log.Debug().Str("chain", chain).Msg("getting report data")

				data, err := rep.GetReportData(
					time.Date(year, 1, 1, 00, 00, 00, 000, time.UTC),
					time.Date(year, 12, 31, 00, 00, 00, 000, time.UTC),
				)
				if err != nil {
					return err
				}

				reports[i] = data

				// Stop the reporter
				rep.Stop()
			}

			var bz []byte
			output, _ := cmd.Flags().GetString(flagOutput)
			switch output {
			case outText:
				bz, err = yaml.Marshal(&reports)
			case outJSON:
				bz, err = json.Marshal(&reports)
			default:
				return fmt.Errorf("invalid output value: %s", output)
			}

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
	cmd.Flags().String(flagOutput, outText, "Type of output (supported values: json, text)")

	return cmd
}