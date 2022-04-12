package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func GetPower(exponent uint64) sdk.Int {
	if exponent == 0 {
		return sdk.OneInt()
	}

	power := sdk.NewInt(10)
	var i uint64 = 0
	for ; i < exponent-1; i++ {
		power = power.Mul(sdk.NewIntFromUint64(10))
	}
	return power
}
