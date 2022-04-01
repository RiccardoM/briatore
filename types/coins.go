package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func GetPower(exponent uint64) sdk.Int {
	power := sdk.NewInt(10)
	var i uint64 = 0
	for ; i < exponent; i++ {
		power = power.Mul(sdk.NewIntFromUint64(10))
	}
	return power
}
