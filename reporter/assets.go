package reporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/riccardom/briatore/types"
)

const (
	assetsListURL = "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"
)

type assetsResponse struct {
	Assets types.Assets `json:"assets"`
}

func (r *Reporter) getBaseNativeDenom(chainName string) (string, error) {
	assets, err := r.getAssetsList()
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

// getAssetsList returns the assets list
func (r *Reporter) getAssetsList() (types.Assets, error) {
	res, err := http.Get(assetsListURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response assetsResponse
	return response.Assets, json.Unmarshal(bz, &response)
}
