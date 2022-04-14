package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocarina/gocsv"

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
	outCSV  = "csv"
)

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

			encodingCfg := simapp.MakeTestEncodingConfig()
			cdc, _ := encodingCfg.Marshaler, encodingCfg.Amino

			var amounts []*types.Amount
			for _, chain := range cfg.Chains {

				// TODO: Get the addresses from the user
				addresses, err := types.GetUniqueSupportedAddresses(chain, cfg.Addresses)
				if err != nil {
					return err
				}

				if len(addresses) == 0 {
					log.Info().Str("chain", chain.Name).Msg("no supported addresses found, skipping")
					continue
				}

				log.Debug().Str("chain", chain.Name).Msg("creating reporter")
				rep, err := reporter.NewReporter(chain, cdc)
				if err != nil {
					log.Error().Str("chain", chain.Name).Err(err).Msg("error while creating the reporter")
					continue
				}

				log.Debug().Str("chain", chain.Name).Msg("getting report data")
				chainAmounts, err := rep.GetAmounts(addresses, date, cfg.Report)
				if err != nil {
					log.Error().Str("chain", chain.Name).Err(err).Msg("error while getting the amounts")
					continue
				}

				amounts = append(amounts, chainAmounts...)
			}

			// Merge the various amounts
			amounts = types.MergeSameAssetsAmounts(amounts)
			amountsOutput := types.Format(amounts)

			var bz []byte
			output, _ := cmd.Flags().GetString(flagOutput)
			switch output {
			case outText:
				bz, err = yaml.Marshal(&amountsOutput)
			case outJSON:
				bz, err = json.Marshal(&amountsOutput)
			case outCSV:
				bz, err = gocsv.MarshalBytes(&amountsOutput)
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
