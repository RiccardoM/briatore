package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"

	"gopkg.in/yaml.v3"

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
		Use:   "report [date] [[chains]]",
		Short: "Reports the data from the given chains for the given date",
		Long: `
Reports the data from the given chains for the given date.

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

			date, err := time.Parse(time.RFC3339, args[0])
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

			encodingCfg := simapp.MakeTestEncodingConfig()
			cdc, _ := encodingCfg.Marshaler, encodingCfg.Amino

			var amounts []types.Amount
			for _, chain := range chains {
				log.Debug().Str("chain", chain).Msg("getting configuration")
				chainCfg := cfg.GetChainConfig(chain)
				if chainCfg == nil {
					return fmt.Errorf("config for chain %s not found", chain)
				}

				log.Debug().Str("chain", chain).Msg("getting account")
				addresses, found := cfg.GetChainAddresses(chain)
				if !found {
					log.Debug().Str("chain", chain).Msg("address not found, skipping")
					continue
				}

				log.Debug().Str("chain", chain).Msg("creating reporter")
				rep, err := reporter.NewReporter(chainCfg, cdc)
				if err != nil {
					return err
				}

				log.Debug().Str("chain", chain).Msg("getting report data")

				chainAmounts, err := rep.GetAmounts(addresses, date, cfg.Report)
				if err != nil {
					return err
				}

				amounts = append(amounts, chainAmounts...)

				// Stop the reporter
				rep.Stop()
			}

			var bz []byte
			output, _ := cmd.Flags().GetString(flagOutput)
			switch output {
			case outText:
				bz, err = yaml.Marshal(&amounts)
			case outJSON:
				bz, err = json.Marshal(&amounts)
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
