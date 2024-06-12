package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

const (
	assetFile     = "assets.json"
	assetsListURL = "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"
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

// --------------------------------------------------------------------------------------------------------------------

type assetsResponse struct {
	Assets Assets `json:"assets"`
}

// GetAssets returns the list of supported assets
func GetAssets() (Assets, error) {
	// Read the stored assets
	bz, err := os.ReadFile(path.Join(HomePath, assetFile))
	if os.IsNotExist(err) {
		// Get the assets from online
		assets, err := RefreshAssets()
		if err != nil {
			return nil, err
		}

		// Return the read assets
		return assets, nil
	}
	if err != nil {
		panic(err)
	}

	var assets Assets
	return assets, json.Unmarshal(bz, &assets)
}

// RefreshAssets gets the assets from the GitHub endpoint and caches them
func RefreshAssets() (Assets, error) {
	res, err := http.Get(assetsListURL)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	bz, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the response
	var response assetsResponse
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return nil, err
	}

	// Store the assets
	err = writeAssets(response.Assets)
	if err != nil {
		return nil, err
	}

	return response.Assets, nil
}

// writeAssets writes the given assets inside the cache
func writeAssets(assets Assets) error {
	bz, err := json.Marshal(&assets)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(HomePath, assetFile), bz, 0600)
}

func GetBaseNativeDenom(chainName string) (string, error) {
	assets, err := GetAssets()
	if err != nil {
		return "", nil
	}

	asset, found := assets.GetAssetByChainName(chainName)
	if !found {
		return "", fmt.Errorf("asset not found")
	}

	nativeDenom, found := asset.GetBaseNativeDenom()
	if !found {
		return "", fmt.Errorf("native denom not found")
	}

	return nativeDenom, nil
}
