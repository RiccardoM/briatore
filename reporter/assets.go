package reporter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/riccardom/briatore/types"
)

const (
	assetsListURL = "https://raw.githubusercontent.com/osmosis-labs/assetlists/main/osmosis-1/osmosis-1.assetlist.json"
)

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

	var assets types.Assets
	return assets, json.Unmarshal(bz, &assets)
}
