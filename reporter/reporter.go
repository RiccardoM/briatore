package reporter

import (
	"context"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/riccardom/briatore/types"
)

type Reporter struct {
	chainCfg  *types.ChainConfig
	reportCfg *types.ReportConfig

	node          *remote.Node
	bankClient    banktypes.QueryClient
	stakingClient stakingtypes.QueryClient
}

func NewReporter(reportCfg *types.ReportConfig, chainCfg *types.ChainConfig, cdc codec.Marshaler) (*Reporter, error) {
	node, err := remote.NewNode(&chainCfg.Node, cdc)
	if err != nil {
		return nil, err
	}

	grpcConnection, err := remote.CreateGrpcConnection(chainCfg.Node.GRPC)
	if err != nil {
		return nil, err
	}

	return &Reporter{
		chainCfg:  chainCfg,
		reportCfg: reportCfg,

		node:          node,
		bankClient:    banktypes.NewQueryClient(grpcConnection),
		stakingClient: stakingtypes.NewQueryClient(grpcConnection),
	}, nil
}

// GetReportData gets the ChainReport
func (r *Reporter) GetReportData(begin, end time.Time) (*types.ChainReport, error) {
	firstReport, err := r.getReport(begin)
	if err != nil {
		return nil, err
	}

	lastReport, err := r.getReport(end)
	if err != nil {
		return nil, err
	}

	return types.NewChaiReport(r.chainCfg.Name, firstReport, lastReport), nil
}

// getReport returns the BalanceReport data for the point in time that is closer to the given timestamp.
// If the provided timestamp is before the genesis, an empty report will be returned instead.
func (r *Reporter) getReport(timestamp time.Time) (*types.BalanceReport, error) {
	block, err := r.getBlockNearTimestamp(timestamp)
	if err != nil {
		return nil, err
	}

	if block == nil {
		// The chain didn't exist at that time, so we just return an empty balance report
		return types.NewBalanceReport(timestamp, r.chainCfg.Address, nil), nil
	}

	return r.getHeightReport(block)
}

// getBlockNearTimestamp returns the block nearest the given timestamp.
// To do this we use the binary search between the genesis height and the latest block time.
func (r *Reporter) getBlockNearTimestamp(timestamp time.Time) (*tmtypes.Block, error) {
	log.Debug().Str("chain", r.chainCfg.Name).Time("timestamp", timestamp).Msg("getting block near timestamp")

	genesis, err := r.node.Genesis()
	if err != nil {
		return nil, fmt.Errorf("error while getting the genesis: %s", err)
	}

	if timestamp.Before(genesis.Genesis.GenesisTime) {
		return nil, nil
	}

	genesisHeight := genesis.Genesis.InitialHeight

	latestHeight, err := r.node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting latest height: %s", err)
	}

	latestBlock, err := r.node.Block(latestHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting latest block: %s", err)
	}

	if timestamp.After(latestBlock.Block.Time) {
		return nil, fmt.Errorf("%s is after latest block time", timestamp)
	}

	// Perform the binary search
	block, err := r.binarySearchBlock(genesisHeight, latestHeight, timestamp)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("chain", r.chainCfg.Name).Time("timestamp", timestamp).Msgf("found block near timestamp: %d", block.Height)

	return block, nil
}

// binarySearchBlock performs a binary search between the given min and max heights,
// searching for the block that is closer to the given timestamp
func (r *Reporter) binarySearchBlock(minHeight, maxHeight int64, timestamp time.Time) (*tmtypes.Block, error) {
	log.Trace().Int64("ming height", minHeight).Int64("max height", maxHeight).Time("timestamp", timestamp).Msg("binary search")

	minBlock, err := r.node.Block(minHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting min block: %s", err)
	}

	if minBlock.Block.Time.Equal(timestamp) {
		// The genesis has the same timestamp of what we are searching for
		return minBlock.Block, nil
	}

	maxBlock, err := r.node.Block(maxHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting max block")
	}

	if maxBlock.Block.Time.Equal(timestamp) {
		// The latest block has the same timestamp we are searching for
		return maxBlock.Block, nil
	}

	if maxBlock.Block.Height-minBlock.Block.Height == 1 {
		// We've reached the point where we only have two blocks.
		// Now we need to find the one that is closer to the searched timestamp
		minDiff := timestamp.Sub(minBlock.Block.Time)
		maxDiff := maxBlock.Block.Time.Sub(timestamp)

		if minDiff < maxDiff {
			// The min block is closer to the timestamp than the max block
			return minBlock.Block, nil
		}

		// The max block is closer to the timestamp than the min block
		return maxBlock.Block, nil
	}

	avgHeight := (maxHeight + minHeight) / 2

	avgBlock, err := r.node.Block(avgHeight)
	if err != nil {
		return nil, fmt.Errorf("error while getting average block: %s", err)
	}

	if avgBlock.Block.Time.Equal(timestamp) {
		// The average block has the same timestamp as the one searched for
		return avgBlock.Block, nil
	}

	if avgBlock.Block.Time.After(timestamp) {
		// If the average block has the timestamp after the searched value, it means the searched
		// value is in between the min value and the average one
		maxHeight = avgBlock.Block.Height
	}

	if avgBlock.Block.Time.Before(timestamp) {
		// If the average block has the timestamp before the searched value, it means the searched
		// value is in between the average block and the max height
		minHeight = avgBlock.Block.Height
	}

	return r.binarySearchBlock(minHeight, maxHeight, timestamp)
}

// getHeightReport returns the BalanceReport for the given height
func (r *Reporter) getHeightReport(block *tmtypes.Block) (*types.BalanceReport, error) {
	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("getting height report")

	ctx := remote.GetHeightRequestContext(context.Background(), block.Height)

	paramsRes, err := r.stakingClient.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	bondDenom := paramsRes.Params.BondDenom

	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("getting balance amount")
	balance, err := r.getBalanceAmount(block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting balance: %s", err)
	}

	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("getting delegations amount")
	delegations, err := r.getDelegationsAmount(block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegations: %s", err)
	}
	balance.Add(delegations...)

	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("getting redelegations amount")
	redelegations, err := r.getReDelegationsAmount(bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while gettig redelegations: %s", err)
	}
	balance.Add(redelegations...)

	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("getting unbonding delegations amount")
	unbondingDelegations, err := r.getUnbondingDelegationsAmount(bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting unbonding delegations: %s", err)
	}
	balance.Add(unbondingDelegations...)

	log.Debug().Str("chain", r.chainCfg.Name).Int64("height", block.Height).Msg("computing report fiat value")
	amount, err := r.getReportAmount(block.Time, balance)
	if err != nil {
		return nil, err
	}

	return types.NewBalanceReport(block.Time, r.chainCfg.Address, amount), nil
}

// getReportAmount returns the corresponding fiat value for the given coins at the provided point in time
func (r *Reporter) getReportAmount(timestamp time.Time, coins sdk.Coins) ([]types.Amount, error) {
	amount := make([]types.Amount, len(coins))
	for i, coin := range coins {
		coingeckoID, found := r.chainCfg.GetCoinGeckoID(coin.Denom)
		if !found {
			amount[i] = types.NewAmount(coin, 0)
			continue
		}

		tokenPrice, err := GetCoinPrice(coingeckoID, timestamp, r.reportCfg.Currency)
		if err != nil {
			return nil, err
		}

		// TODO: This cast might not be safe
		tokenValue := tokenPrice * float64(coin.Amount.Int64())
		amount[i] = types.NewAmount(coin, tokenValue)
	}
	return amount, nil
}

func (r *Reporter) Stop() {
	r.node.Stop()
}
