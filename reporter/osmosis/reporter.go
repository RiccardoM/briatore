package osmosis

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	gammtypes "github.com/osmosis-labs/osmosis/v25/x/gamm/types"
	lockuptypes "github.com/osmosis-labs/osmosis/v25/x/lockup/types"
	poolmanagergrpc "github.com/osmosis-labs/osmosis/v25/x/poolmanager/client/queryproto"

	"github.com/riccardom/briatore/types"
	"github.com/riccardom/briatore/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Reporter struct {
	cdc codec.Codec

	grpcHeaders map[string]string

	gammQueryClient        gammtypes.QueryClient
	poolmanagerQueryClient poolmanagergrpc.QueryClient
	lockupQueryClient      lockuptypes.QueryClient
}

func NewReporter(grpcConnection grpc.ClientConnInterface, grpcHeaders map[string]string, cdc codec.Codec) (*Reporter, error) {
	return &Reporter{
		cdc:                    cdc,
		grpcHeaders:            grpcHeaders,
		gammQueryClient:        gammtypes.NewQueryClient(grpcConnection),
		poolmanagerQueryClient: poolmanagergrpc.NewQueryClient(grpcConnection),
		lockupQueryClient:      lockuptypes.NewQueryClient(grpcConnection),
	}, nil
}

func (r *Reporter) GetAmount(address string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting amount")

	gammBalance := sdk.NewCoins()

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting unlockable amount")
	unlockableAmount, err := r.getUnlockableAmount(address, height)
	if err != nil {
		return nil, err
	}
	gammBalance = gammBalance.Add(unlockableAmount...)

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting unlocking amount")
	unlockingAmount, err := r.getUnlockingAmount(address, height)
	if err != nil {
		return nil, err
	}
	gammBalance = gammBalance.Add(unlockingAmount...)

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting locked amount")
	lockedAmount, err := r.getLockedAmount(address, height)
	if err != nil {
		return nil, err
	}
	gammBalance = gammBalance.Add(lockedAmount...)

	var balance sdk.Coins
	for _, gammToken := range gammBalance {
		amount, err := r.convertPoolShares(gammToken, height)
		if err != nil {
			return nil, err
		}
		balance = balance.Add(amount...)
	}

	return balance, nil
}

// convertPoolShares converts the given GAMM token into the proper denoms
func (r *Reporter) convertPoolShares(gammToken sdk.Coin, height int64) (sdk.Coins, error) {
	ctx := utils.GetRequestContext(height, r.grpcHeaders)

	// Get the pool id
	poolID, err := utils.ParsePoolID(gammToken.Denom)
	if err != nil {
		return nil, err
	}

	// Get the pool liquidity and total shares
	poolLiquidityRes, err := r.poolmanagerQueryClient.TotalPoolLiquidity(ctx, &poolmanagergrpc.TotalPoolLiquidityRequest{PoolId: poolID})
	if err != nil {
		return nil, fmt.Errorf("error while querying the pool: %w", err)
	}
	poolLiquidity := poolLiquidityRes.Liquidity

	poolSharesRes, err := r.gammQueryClient.TotalShares(ctx, &gammtypes.QueryTotalSharesRequest{PoolId: poolID})
	if err != nil {
		return nil, fmt.Errorf("error while querying the pool: %w", err)
	}
	poolShares := poolSharesRes.TotalShares

	// Compute the share ratio
	shareRatio := gammToken.Amount.ToLegacyDec().QuoInt(poolShares.Amount).MulInt(types.GetPower(2))

	balance := sdk.NewCoins()
	for _, asset := range poolLiquidity {
		balance = balance.Add(sdk.NewCoin(asset.Denom, shareRatio.MulInt(asset.Amount).RoundInt().Quo(sdk.NewInt(100))))
	}

	return balance, nil
}

func (r *Reporter) getUnlockableAmount(address string, height int64) (sdk.Coins, error) {
	ctx := utils.GetRequestContext(height, r.grpcHeaders)
	res, err := r.lockupQueryClient.AccountUnlockableCoins(ctx, &lockuptypes.AccountUnlockableCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}

func (r *Reporter) getUnlockingAmount(address string, height int64) (sdk.Coins, error) {
	ctx := utils.GetRequestContext(height, r.grpcHeaders)
	res, err := r.lockupQueryClient.AccountUnlockingCoins(ctx, &lockuptypes.AccountUnlockingCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}

func (r *Reporter) getLockedAmount(address string, height int64) (sdk.Coins, error) {
	ctx := utils.GetRequestContext(height, r.grpcHeaders)
	res, err := r.lockupQueryClient.AccountLockedCoins(ctx, &lockuptypes.AccountLockedCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}
