package types

import (
	sdkmath "cosmossdk.io/math"
)

func GetPower(exponent uint64) sdkmath.Int {
	if exponent == 0 {
		return sdkmath.OneInt()
	}

	power := sdkmath.NewInt(10)
	var i uint64 = 0
	for ; i < exponent-1; i++ {
		power = power.Mul(sdkmath.NewIntFromUint64(10))
	}
	return power
}
