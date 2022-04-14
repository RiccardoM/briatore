package types

import "github.com/cosmos/cosmos-sdk/types/bech32"

// GetUniqueSupportedAddresses returns the list of all the given addresses that are supported by the
// provided chain config, removing any duplicated address that might be specified for different chains
func GetUniqueSupportedAddresses(chainCfg *ChainConfig, addresses []string) ([]string, error) {
	var supportedAddresses = map[string]int{}
	for _, address := range addresses {
		prefix, _, err := bech32.DecodeAndConvert(address)
		if err != nil {
			return nil, err
		}

		if chainCfg.Bech32Prefix == prefix {
			if _, exists := supportedAddresses[address]; !exists {
				supportedAddresses[address] = 1
			}
		}
	}

	var slice = make([]string, len(supportedAddresses))
	var i = 0
	for address := range supportedAddresses {
		slice[i] = address
		i++
	}

	return slice, nil
}
