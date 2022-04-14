package reporter

import (
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/remote"
	"github.com/rs/zerolog/log"

	"github.com/riccardom/briatore/types"
)

type Reporter struct {
	cdc codec.Codec

	chain *types.ChainConfig

	grpcConnection *grpc.ClientConn
	grpcHeaders    map[string]string

	node *remote.Node

	bankClient    banktypes.QueryClient
	stakingClient stakingtypes.QueryClient
}

func NewReporter(cfg *types.ChainConfig, cdc codec.Codec) (*Reporter, error) {
	grpcAddress, insecure, headers := types.ParseGRPCAddress(cfg.GRPCAddress)
	remoteCfg := remote.NewDetails(
		remote.NewRPCConfig("briatore", cfg.RPCAddress, 10),
		remote.NewGrpcConfig(grpcAddress, insecure),
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
		cdc:            cdc,
		chain:          cfg,
		grpcConnection: grpcConnection,
		grpcHeaders:    headers,
		node:           node,
		bankClient:     banktypes.NewQueryClient(grpcConnection),
		stakingClient:  stakingtypes.NewQueryClient(grpcConnection),
	}, nil
}

// GetAmounts returns the amount that the given addresses hold at the point in time that is closest to the given timestamp.
// If the provided timestamp is before the genesis, an empty report will be returned instead.
// NOTE. Calling this method will close the node as soon as it returns
func (r *Reporter) GetAmounts(addresses []string, timestamp time.Time, cfg *types.ReportConfig) ([]*types.Amount, error) {
	defer r.stop()

	blockData, err := r.getBlockNearTimestamp(timestamp)
	if err != nil {
		return nil, err
	}

	// Get the overall hold amount
	sum := sdk.NewCoins()
	for _, address := range addresses {
		amount, err := r.getHeightAmount(address, blockData.Height)
		if err != nil {
			return nil, err
		}
		sum = sum.Add(amount...)
	}

	// Get the amounts
	return r.getCoinsAmounts(blockData.Timestamp, sum, cfg)
}

// getHeightAmount returns the hold amount at the given height
func (r *Reporter) getHeightAmount(address string, height int64) (sdk.Coins, error) {
	if height == 0 {
		// If the height is 0 it means the chain didn't exist, so we just return an empty amount
		return nil, nil
	}

	log.Debug().Str("chain", r.chain.Name).Str("address", address).Int64("height", height).Msg("getting height report")

	bondDenom, err := types.GetBaseNativeDenom(r.chain.Name)
	if err != nil {
		return nil, fmt.Errorf("error while getting base native denom: %s", err)
	}

	balance, err := r.getBalanceAmount(address, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting balance: %s", err)
	}

	delegations, err := r.getDelegationsAmount(address, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegations: %s", err)
	}
	balance = balance.Add(delegations...)

	redelegations, err := r.getReDelegationsAmount(address, bondDenom, height)
	if err != nil {
		return nil, fmt.Errorf("error while gettig redelegations: %s", err)
	}
	balance = balance.Add(redelegations...)

	unbondingDelegations, err := r.getUnbondingDelegationsAmount(address, bondDenom, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting unbonding delegations: %s", err)
	}
	balance = balance.Add(unbondingDelegations...)

	osmosisAmount, err := r.getOsmosisAmount(address, height)
	if err != nil {
		return nil, fmt.Errorf("error while getting osmosis amount: %s", err)
	}
	balance = balance.Add(osmosisAmount...)

	return balance, nil
}

// getCoinsAmounts returns the corresponding fiat value for the given coins at the provided point in time
func (r *Reporter) getCoinsAmounts(timestamp time.Time, coins sdk.Coins, cfg *types.ReportConfig) ([]*types.Amount, error) {
	log.Debug().Str("chain", r.chain.Name).Time("timestamp", timestamp).Msg("computing report fiat value")

	assets, err := types.GetAssets()
	if err != nil {
		return nil, err
	}

	var amounts []*types.Amount
	for _, coin := range coins {
		// Get the CoinGecko ID, if not found just return a value of 0
		asset, found := assets.GetAssetByCoinDenom(coin.Denom)
		if !found {
			log.Info().Str("denom", coin.Denom).Msg("asset not found")
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
