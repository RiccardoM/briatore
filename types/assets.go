package types

import (
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Asset struct {
	Name        string                 `json:"name"`
	Base        string                 `json:"base"`
	Symbol      string                 `json:"symbol"`
	CoingeckoID string                 `json:"coingecko_id"`
	DenomUnits  []*banktypes.DenomUnit `json:"denom_units"`
}

func (a *Asset) GetBaseNativeDenom() (nativeDenom string, found bool) {
	for _, unit := range a.DenomUnits {
		if unit.Denom == a.Base {
			if len(unit.Aliases) > 0 {
				return unit.Aliases[0], true
			}
			return unit.Denom, true
		}
	}
	return "", false
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

func (l Assets) GetAssetByChainName(chainName string) (asset *Asset, found bool) {
	for _, asset := range l {
		if strings.EqualFold(asset.Name, chainName) {
			return asset, true
		}
	}
	return nil, false
}

func (l Assets) GetAssetByCoinDenom(coinDenom string) (asset *Asset, found bool) {
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
