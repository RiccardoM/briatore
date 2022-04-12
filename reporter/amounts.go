package reporter

import (
	"context"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v3/node/remote"
)

func (r *Reporter) getBalanceAmount(address string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", height).Msg("getting balance amount")

	ctx := remote.GetHeightRequestContext(context.Background(), height)

	balance := sdk.NewCoins()
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := r.bankClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{
			Address: address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		balance = balance.Add(res.Balances...)
		nextKey = res.Pagination.NextKey
		stop = len(nextKey) == 0
	}

	return balance, nil
}

func (r *Reporter) getDelegationsAmount(address string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", height).Msg("getting delegations amount")

	ctx := remote.GetHeightRequestContext(context.Background(), height)

	var delegations []stakingtypes.DelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := r.stakingClient.DelegatorDelegations(ctx, &stakingtypes.QueryDelegatorDelegationsRequest{
			DelegatorAddr: address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, res.DelegationResponses...)
		nextKey = res.Pagination.NextKey
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		amount = amount.Add(res.Balance)
	}

	return amount, nil
}

func (r *Reporter) getReDelegationsAmount(address string, bondDenom string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", height).Msg("getting redelegations amount")

	ctx := remote.GetHeightRequestContext(context.Background(), height)
	var delegations []stakingtypes.RedelegationResponse
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := r.stakingClient.Redelegations(ctx, &stakingtypes.QueryRedelegationsRequest{
			DelegatorAddr: address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, res.RedelegationResponses...)
		nextKey = res.Pagination.NextKey
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		for _, entry := range res.Entries {
			amount = amount.Add(sdk.NewCoin(bondDenom, entry.Balance))
		}
	}

	return amount, nil
}

func (r *Reporter) getUnbondingDelegationsAmount(address string, bondDenom string, height int64) (sdk.Coins, error) {
	log.Debug().Str("chain", r.cfg.Name).Int64("height", height).Msg("getting unbonding delegations amount")

	ctx := remote.GetHeightRequestContext(context.Background(), height)
	var delegations []stakingtypes.UnbondingDelegation
	var nextKey []byte
	var stop = false
	for !stop {
		res, err := r.stakingClient.DelegatorUnbondingDelegations(ctx, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
			DelegatorAddr: address,
			Pagination: &query.PageRequest{
				Key: nextKey,
			},
		})
		if err != nil {
			return nil, err
		}

		delegations = append(delegations, res.UnbondingResponses...)
		nextKey = res.Pagination.NextKey
		stop = len(nextKey) == 0
	}

	amount := sdk.NewCoins()
	for _, res := range delegations {
		for _, entry := range res.Entries {
			amount = amount.Add(sdk.NewCoin(bondDenom, entry.Balance))
		}
	}

	return amount, nil
}
