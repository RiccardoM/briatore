package reporter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"

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
	cfg *types.ChainConfig

	grpcConnection *grpc.ClientConn
	node           *remote.Node

	bankClient    banktypes.QueryClient
	stakingClient stakingtypes.QueryClient
}

func NewReporter(cfg *types.ChainConfig, cdc codec.Codec) (*Reporter, error) {
	remoteCfg := remote.NewDetails(
		remote.NewRPCConfig("briatore", cfg.RPCAddress, 10),
		remote.NewGrpcConfig(cfg.GRPCAddress, strings.Contains(cfg.GRPCAddress, "https://")),
	)
	node, err := remote.NewNode(remoteCfg, cdc)
	if err != nil {
		return nil, err
	}

	grpcConnection, err := remote.CreateGrpcConnection(remoteCfg.GRPC)
	if err != nil {
		return nil, err
	}

	return &Reporter{
		cfg:            cfg,
		grpcConnection: grpcConnection,
		node:           node,
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		stakingClient:  stakingtypes.NewQueryClient(grpcConnection),
	}, nil
}

// GetReports returns the BalanceReport data for the point in time that is closer to the given timestamp.
// If the provided timestamp is before the genesis, an empty report will be returned instead.
func (r *Reporter) GetReports(addresses []string, timestamp time.Time, cfg *types.ReportConfig) (types.BalancesReports, error) {
	block, err := r.getBlockNearTimestamp(timestamp)
	if err != nil {
		return nil, err
	}

	if block == nil {
		// The chain didn't exist at that time, so we just return an empty balance report
		return nil, nil
	}

	return r.getHeightReports(addresses, block, cfg)
}

// getHeightReports returns the list of BalanceReports for the given height
func (r *Reporter) getHeightReports(addresses []string, block *tmtypes.Block, cfg *types.ReportConfig) (types.BalancesReports, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", block.Height).Msg("getting height reports")

	reports := make(types.BalancesReports, len(addresses))
	for i, address := range addresses {
		report, err := r.getHeightReport(address, block, cfg)
		if err != nil {
			return nil, err
		}
		reports[i] = report
	}

	return reports, nil
}

// getHeightReport returns the BalanceReport for the given height
func (r *Reporter) getHeightReport(address string, block *tmtypes.Block, cfg *types.ReportConfig) (*types.BalanceReport, error) {
	log.Debug().Str("chain", r.cfg.Name).Str("address", address).Int64("height", block.Height).Msg("getting height report")

	ctx := remote.GetHeightRequestContext(context.Background(), block.Height)

	paramsRes, err := r.stakingClient.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}
	bondDenom := paramsRes.Params.BondDenom

	balance, err := r.getBalanceAmount(address, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting balance: %s", err)
	}

	delegations, err := r.getDelegationsAmount(address, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegations: %s", err)
	}
	balance.Add(delegations...)

	redelegations, err := r.getReDelegationsAmount(address, bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while gettig redelegations: %s", err)
	}
	balance.Add(redelegations...)

	unbondingDelegations, err := r.getUnbondingDelegationsAmount(address, bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting unbonding delegations: %s", err)
	}
	balance.Add(unbondingDelegations...)

	osmosisAmount, err := r.getOsmosisAmount(address, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting osmosis amount: %s", err)
	}
	balance.Add(osmosisAmount...)

	amount, err := r.getReportAmount(block.Time, balance, cfg)
	if err != nil {
		return nil, err
	}

	return types.NewBalanceReport(block.Time, address, amount), nil
}

// getReportAmount returns the corresponding fiat value for the given coins at the provided point in time
func (r *Reporter) getReportAmount(timestamp time.Time, coins sdk.Coins, cfg *types.ReportConfig) ([]types.Amount, error) {
	log.Debug().Str("chain", r.cfg.Name).Time("timestamp", timestamp).Msg("computing report fiat value")

	assets, err := r.getAssetsList()
	if err != nil {
		return nil, err
	}

	amount := make([]types.Amount, len(coins))
	for i, coin := range coins {
		// Get the CoinGecko ID, if not found just return a value of 0
		asset, found := assets.GetAsset(coin.Denom)
		if !found {
			amount[i] = types.NewAmount(coin, "0")
			continue
		}

		// Get the token price
		tokenPrice, err := GetCoinPrice(asset.CoingeckoID, timestamp, cfg.Currency)
		if err != nil {
			return nil, err
		}
		tokenPriceDec, err := sdk.NewDecFromStr(fmt.Sprintf("%.2f", tokenPrice))
		if err != nil {
			return nil, err
		}

		// Compute the token value
		tokenAmount := coin.Amount.Quo(types.GetPower(asset.GetMaxExponent()))
		tokenValue := tokenAmount.ToDec().Mul(tokenPriceDec)

		amount[i] = types.NewAmount(coin, tokenValue.String())
	}

	return amount, nil
}

func (r *Reporter) Stop() {
	r.node.Stop()
}
