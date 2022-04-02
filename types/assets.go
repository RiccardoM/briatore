package types

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Asset struct {
	Name        string                 `json:"name"`
	CoingeckoID string                 `json:"coingecko_id"`
	DenomUnits  []*banktypes.DenomUnit `json:"denom_units"`
}

func (a *Asset) GetMaxExponent() uint64 {
	var maxExponent uint64 = 0
	for _, unit := range a.DenomUnits {
		if uint64(unit.Exponent) > maxExponent {
			maxExponent = uint64(unit.Exponent)
		}
	}
	return maxExponent
}

// --------------------------------------------------------------------------------------------------------------------

type Assets []*Asset

func (l Assets) GetAsset(coinDenom string) (asset *Asset, found bool) {
	for _, asset := range l {
		for _, denom := range asset.DenomUnits {
			if denom.Denom == coinDenom {
				return asset, true
			}
			for _, alias := range denom.Aliases {
				if alias == coinDenom {
					return asset, true
				}
			}
		}
	}
	return nil, false
}
