package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Amount struct {
	Asset  *Asset  `yaml:"asset" json:"asset"`
	Amount sdk.Dec `yaml:"amount" json:"amount"`
	Value  sdk.Dec `yaml:"value" json:"value"`
}

func NewAmount(asset *Asset, amount sdk.Dec, value sdk.Dec) *Amount {
	return &Amount{
		Asset:  asset,
		Amount: amount,
		Value:  value,
	}
}

// --------------------------------------------------------------------------------------------------------------------
// CSV Support

type AmountOutput struct {
	Asset  string `json:"asset" yaml:"asset" csv:"asset"`
	Amount string `json:"amount" yaml:"amount" csv:"amount"`
	Value  string `json:"value" yaml:"value" csv:"value"`
}

// Format formats the given amounts to be later printed properly
func Format(amounts []*Amount) []AmountOutput {
	csvAmounts := make([]AmountOutput, len(amounts))
	for i, amount := range amounts {
		csvAmounts[i] = AmountOutput{
			Asset:  amount.Asset.Symbol,
			Amount: amount.Amount.String(),
			Value:  amount.Value.String(),
		}
	}
	return csvAmounts
}

// --------------------------------------------------------------------------------------------------------------------

// MergeSameAssetsAmounts merges together the various amounts for the same assets present inside the given slice
func MergeSameAssetsAmounts(slice []*Amount) []*Amount {
	assets := map[string]*Asset{}
	amounts := map[string]sdk.Dec{}
	values := map[string]sdk.Dec{}

	// Collect all the unique assets
	for _, amount := range slice {

		// Store the asset
		if _, ok := assets[amount.Asset.Name]; !ok {
			assets[amount.Asset.Name] = amount.Asset
		}

		// Store the amounts
		assetAmount, ok := amounts[amount.Asset.Name]
		if !ok {
			assetAmount = sdk.ZeroDec()
		}
		amounts[amount.Asset.Name] = assetAmount.Add(amount.Amount)

		// Store the values
		assetValue, ok := values[amount.Asset.Name]
		if !ok {
			assetValue = sdk.ZeroDec()
		}
		values[amount.Asset.Name] = assetValue.Add(amount.Value)
	}

	var result []*Amount
	for name, asset := range assets {
		result = append(result, NewAmount(asset, amounts[name], values[name]))
	}

	return result
}
