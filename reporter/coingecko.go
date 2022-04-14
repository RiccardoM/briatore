package reporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/riccardom/briatore/types"
)

const (
	CoingeckoEndpoint = "https://api.coingecko.com/api/v3/coins/{id}/history?date={date}"
)

// GetCoinPrice gets the historical price of the coin having the given CoinGecko ID,
// measured in the given currency.
func GetCoinPrice(id string, timestamp time.Time, currency string) (float64, error) {
	priceData, found, err := types.GetPriceData(id, currency, timestamp)
	if err != nil {
		return 0, err
	}

	if !found {
		price, err := getPriceFromAPI(id, timestamp, currency)
		if err != nil {
			return 0, err
		}

		priceData = types.NewPriceData(id, price, currency, timestamp)

		// Cache the price data
		err = types.CachePriceData(priceData)
		if err != nil {
			return 0, err
		}
	}

	return priceData.Price, nil
}

// getPriceFromAPI returns the price for the coin having the given id for the given timestamp and currency
func getPriceFromAPI(id string, timestamp time.Time, currency string) (float64, error) {
	log.Debug().Str("id", id).Time("timestamp", timestamp).Msg("getting price from API")

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

	return response.GetCoinPrice(currency)
}
