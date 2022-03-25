package report

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/riccardom/briatore/reporter"
	"github.com/riccardom/briatore/types"
)

func GetReportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "report [year] [[chains]]",
		Short: "Reports the data from the given chains for the given year",
		Long: `
Reports the data from the given chains for the given year.

If no chain is provided, then all chains present inside the configuration file will be reported.
Multiple chains can be specified separating them using spaces.`,
		Example: "report cosmos-hub osmosis chihuahua",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			for _, chain := range chains {
				log.Debug().Str("chain", chain).Msg("getting configuration")

				chainCfg, found := cfg.GetChainConfig(chain)
				if !found {
					return fmt.Errorf("config for chain %s not found", chain)
				}

				rep, err := reporter.NewReporter(chainCfg, cdc)
				if err != nil {
					return err
				}

				log.Debug().Str("chain", chain).Msg("getting report data")

				data, err := rep.GetReportData(
					time.Date(year, 1, 1, 00, 00, 00, 000, time.UTC),
					time.Date(year, 12, 31, 00, 00, 00, 000, time.UTC),
					cfg.Report, // TODO: Allow to customize this with flags
				)

				// TODO: Output this somewhere
				print(data)
			}

			return nil
		},
	}
}
