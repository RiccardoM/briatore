package osmosis

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v3/node/remote"
	lockuptypes "github.com/osmosis-labs/osmosis/v7/x/lockup/types"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type OsmosisReporter struct {
	lockupQueryClient lockuptypes.QueryClient
}

func NewOsmosisReporter(grpcConnection *grpc.ClientConn) (*OsmosisReporter, error) {
	return &OsmosisReporter{
		lockupQueryClient: lockuptypes.NewQueryClient(grpcConnection),
	}, nil
}

func (r *OsmosisReporter) GetAmount(address string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting amount")

	balance := sdk.NewCoins()

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting unlockable amount")
	unlockableAmount, err := r.getUnlockableAmount(address, height)
	if err != nil {
		return nil, err
	}
	balance.Add(unlockableAmount...)

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting unlocking amount")
	unlockingAmount, err := r.getUnlockingAmount(address, height)
	if err != nil {
		return nil, err
	}
	balance.Add(unlockingAmount...)

	log.Debug().Str("chain", "osmosis").Int64("height", height).Msg("getting locked amount")
	lockedAmount, err := r.getLockedAmount(address, height)
	if err != nil {
		return nil, err
	}
	balance.Add(lockedAmount...)

	return balance, nil
}

func (r *OsmosisReporter) getUnlockableAmount(address string, height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)
	res, err := r.lockupQueryClient.AccountUnlockableCoins(ctx, &lockuptypes.AccountUnlockableCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}

func (r *OsmosisReporter) getUnlockingAmount(address string, height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)
	res, err := r.lockupQueryClient.AccountUnlockingCoins(ctx, &lockuptypes.AccountUnlockingCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}

func (r *OsmosisReporter) getLockedAmount(address string, height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)
	res, err := r.lockupQueryClient.AccountLockedCoins(ctx, &lockuptypes.AccountLockedCoinsRequest{
		Owner: address,
	})
	if err != nil {
		return nil, err
	}

	return res.Coins, nil
}
