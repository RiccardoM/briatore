package report

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/osmosis-labs/osmosis/v7/app"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/riccardom/briatore/reporter"
	"github.com/riccardom/briatore/types"
)

// GetReport returns the serialized report bytes for the given configuration, addresses and date.
// The report will be serialized properly based on the given output type.
func GetReport(cfg *types.Config, addresses []string, date time.Time) *types.ReportResult {
	cdc, _ := app.MakeCodecs()

	var amounts []*types.Amount
	for _, chain := range cfg.Chains {
		log.Info().Str("chain", chain.Name).Msg("getting report")

		addresses, err := types.GetUniqueSupportedAddresses(chain, addresses)
		if err != nil {
			return types.NewErrorReportResult(err)
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

		log.Info().Str("chain", chain.Name).Msg("report retrieved")
	}

	// Merge the various amounts and format them
	return types.NewAmountsReportResult(types.Format(types.MergeSameAssetsAmounts(amounts)))
}

// MarshalAmounts marshals the given amount based on the provided output
func MarshalAmounts(amounts []types.AmountOutput, output types.Output) ([]byte, error) {
	switch output {
	case types.OutText:
		return yaml.Marshal(&amounts)
	case types.OutJSON:
		return json.Marshal(&amounts)
	case types.OutCSV:
		return gocsv.MarshalBytes(&amounts)
	default:
		return nil, fmt.Errorf("invalid output value: %s", output)
	}
}
