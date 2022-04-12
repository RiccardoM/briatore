package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Amount struct {
	Asset  *Asset  `yaml:"asset" json:"asset"`
	Amount sdk.Dec `yaml:"amount" json:"amount"`
	Value  sdk.Dec `yaml:"value" json:"value"`
}

func NewAmount(asset *Asset, amount sdk.Dec, value sdk.Dec) Amount {
	return Amount{
		Asset:  asset,
		Amount: amount,
		Value:  value,
	}
}
