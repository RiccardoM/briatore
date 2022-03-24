package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/riccardom/briatore/types"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	CoingeckoEndpoint = "https://api.coingecko.com/api/v3/coins/{id}/history?date={date}"
)

// GetCoinPrice gets the historical price of the coin having the given CoinGecko ID,
// measured in the given currency.
func GetCoinPrice(id string, timestamp time.Time, currency string) (float64, error) {
	endpoint := strings.ReplaceAll(CoingeckoEndpoint, "{id}", id)
	endpoint = strings.ReplaceAll(endpoint, "{date}", timestamp.Format("02-01-2006"))

	res, err := http.Get(endpoint)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad token price history response: status %d", res.StatusCode)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var response types.HistoryResponse
	err = json.Unmarshal(bz, &response)
	if err != nil {
		return 0, err
	}

	price, err := response.GetCoinPrice(currency)
	if err != nil {
		return 0, err
	}

	return price, nil
}
