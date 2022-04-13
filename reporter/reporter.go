package reporter

import (
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
		remote.NewGrpcConfig(cfg.GRPCAddress, !strings.Contains(cfg.GRPCAddress, "https://")),
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

// GetAmounts returns the amounts for the point in time that is closer to the given timestamp.
// If the provided timestamp is before the genesis, an empty report will be returned instead.
// NOTE. Calling this method will close the node as soon as it returns
func (r *Reporter) GetAmounts(addresses []string, timestamp time.Time, cfg *types.ReportConfig) ([]types.Amount, error) {
	defer r.stop()

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
func (r *Reporter) getHeightReports(addresses []string, block *tmtypes.Block, cfg *types.ReportConfig) ([]types.Amount, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", block.Height).Msg("getting height reports")

	var result []types.Amount
	for _, address := range addresses {
		amounts, err := r.getHeightReport(address, block, cfg)
		if err != nil {
			return nil, err
		}

		result = append(result, amounts...)
	}

	return result, nil
}

// getHeightReport returns the BalanceReport for the given height
func (r *Reporter) getHeightReport(address string, block *tmtypes.Block, cfg *types.ReportConfig) ([]types.Amount, error) {
	log.Debug().Str("chain", r.cfg.Name).Str("address", address).Int64("height", block.Height).Msg("getting height report")

	bondDenom, err := r.getBaseNativeDenom(r.cfg.Name)
	if err != nil {
		return nil, fmt.Errorf("error while getting base native denom: %s", err)
	}

	balance, err := r.getBalanceAmount(address, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting balance: %s", err)
	}

	delegations, err := r.getDelegationsAmount(address, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegations: %s", err)
	}
	balance = balance.Add(delegations...)

	redelegations, err := r.getReDelegationsAmount(address, bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while gettig redelegations: %s", err)
	}
	balance = balance.Add(redelegations...)

	unbondingDelegations, err := r.getUnbondingDelegationsAmount(address, bondDenom, block.Height)
	if err != nil {
		return nil, fmt.Errorf("error while getting unbonding delegations: %s", err)
	}
	balance = balance.Add(unbondingDelegations...)

	return r.getReportAmount(block.Time, balance, cfg)
}

// getReportAmount returns the corresponding fiat value for the given coins at the provided point in time
func (r *Reporter) getReportAmount(timestamp time.Time, coins sdk.Coins, cfg *types.ReportConfig) ([]types.Amount, error) {
	log.Debug().Str("chain", r.cfg.Name).Time("timestamp", timestamp).Msg("computing report fiat value")

	assets, err := r.getAssetsList()
	if err != nil {
		return nil, err
	}

	var amounts []types.Amount
	for _, coin := range coins {
		// Get the CoinGecko ID, if not found just return a value of 0
		asset, found := assets.GetAssetByCoinDenom(coin.Denom)
		if !found {
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
		tokenAmount := coin.Amount.ToDec().QuoInt(types.GetPower(asset.GetMaxExponent()))
		tokenValue := tokenAmount.Mul(tokenPriceDec)

		amounts = append(amounts, types.NewAmount(asset, tokenAmount, tokenValue))
	}

	return amounts, nil
}

func (r *Reporter) stop() {
	r.node.Stop()
}
