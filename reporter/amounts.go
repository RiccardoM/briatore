package reporter

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/remote"
)

func (r *Reporter) getBalanceAmount(height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)

	balance := sdk.NewCoins()
	var nextKey []byte
	var stop = false
	for !stop {
		balRes, err := r.bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
			Address: r.chainCfg.Address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		balance.Add(balRes.Balances...)
		stop = len(nextKey) == 0
	}

	return balance, nil
}

func (r *Reporter) getDelegationsAmount(height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		delRes, err := r.stakingClient.DelegatorDelegations(ctx, &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: r.chainCfg.Address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, delRes.DelegationResponses...)
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		amount.Add(res.Balance)
	}

	return amount, nil
}

func (r *Reporter) getReDelegationsAmount(bondDenom string, height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)

	var delegations []stakingtypes.RedelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		delRes, err := r.stakingClient.Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{
			DelegatorAddr: r.chainCfg.Address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, delRes.RedelegationResponses...)
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		for _, entry := range res.Entries {
			amount.Add(sdk.NewCoin(bondDenom, entry.Balance))
		}
	}

	return amount, nil
}

func (r *Reporter) getUnbondingDelegationsAmount(bondDenom string, height int64) (sdk.Coins, error) {
	ctx := remote.GetHeightRequestContext(context.Background(), height)

	var delegations []stakingtypes.UnbondingDelegation
	var nextKey []byte
	var stop = false
	for !stop {
		delRes, err := r.stakingClient.DelegatorUnbondingDelegations(ctx, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: r.chainCfg.Address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, delRes.UnbondingResponses...)
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		for _, entry := range res.Entries {
			amount.Add(sdk.NewCoin(bondDenom, entry.Balance))
		}
	}

	return amount, nil
}
